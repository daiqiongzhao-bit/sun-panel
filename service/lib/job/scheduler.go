package job

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sun-panel/global"
	"sun-panel/models"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	scheduler *cron.Cron
	jobEntryMap sync.Map // jobTaskId -> cron.EntryID
)

// InitScheduler 初始化定时任务调度器，启动时从数据库加载所有运行中的任务
func InitScheduler() {
	scheduler = cron.New(cron.WithSeconds())
	loadJobsFromDB()
	scheduler.Start()
	global.Logger.Info("定时任务调度器已启动")
}

// loadJobsFromDB 从数据库加载运行中的任务并注册
func loadJobsFromDB() {
	jobModel := models.JobTask{}
	jobs, err := jobModel.GetAll()
	if err != nil {
		global.Logger.Error("加载定时任务失败: " + err.Error())
		return
	}

	for _, job := range jobs {
		addJobToScheduler(job)
	}
}

// addJobToScheduler 将任务添加到调度器
func addJobToScheduler(job models.JobTask) {
	entryId, err := scheduler.AddFunc(job.CronExpr, func() {
		executeJob(job.ID, job.Name, job.Content, job.JobType)
	})
	if err != nil {
		global.Logger.Error(fmt.Sprintf("注册定时任务失败[%s] cron=%s: %s", job.Name, job.CronExpr, err.Error()))
		return
	}
	jobEntryMap.Store(job.ID, entryId)

	// 计算下次执行时间
	entry := scheduler.Entry(entryId)
	nextRun := entry.Next
	(&models.JobTask{}).UpdateRunTime(job.ID, nil, &nextRun)
}

// removeJobFromScheduler 从调度器移除任务
func removeJobFromScheduler(jobId uint) {
	if entryId, ok := jobEntryMap.Load(jobId); ok {
		scheduler.Remove(entryId.(cron.EntryID))
		jobEntryMap.Delete(jobId)
	}
}

// RegisterJob 注册任务到调度器
func RegisterJob(job models.JobTask) error {
	addJobToScheduler(job)
	return nil
}

// UnregisterJob 从调度器移除任务
func UnregisterJob(jobId uint) {
	removeJobFromScheduler(jobId)
}

// executeJob 执行单个任务（核心逻辑）
func executeJob(jobId uint, jobName, content string, jobType int) {
	startTime := time.Now()
	log := models.JobLog{
		JobId:     jobId,
		JobName:   jobName,
		StartTime: startTime.Format("2006-01-02 15:04:05"),
		Status:    models.JOB_STATUS_RUNNING,
	}

	defer func() {
		if r := recover(); r != nil {
			log.Duration = time.Since(startTime).Milliseconds()
			log.Status = 2
			errMsg := fmt.Sprintf("panic: %v", r)
			log.ErrorMsg = errMsg
			log.Create()
			// panic 时也通知超管
			notifyAdminOnFailure(jobId, jobName, errMsg)
		}
	}()

	var execErr error

	switch jobType {
	case JOB_TYPE_REMINDER:
		// 定时提醒：创建一条通知
		notice := models.Notice{
			Title:       fmt.Sprintf("[定时提醒] %s", jobName),
			Content:     content,
			DisplayType: models.NOTICE_DISPLAY_TYPE_HOME,
			NoticeType:  models.NOTICE_TYPE_ANNOUNCEMENT,
			OneRead:     1,
			Status:      models.NOTICE_STATUS_ENABLED,
			UserId:      0, // 系统触发
		}
		if err := global.Db.Create(&notice).Error; err != nil {
			execErr = err
		}

	case JOB_TYPE_HEALTH_CHECK:
		// 导航链接健康检查
		execErr = performHealthCheck()

	default:
		execErr = fmt.Errorf("未知的任务类型: %d", jobType)
	}

	// 记录执行日志
	log.Duration = time.Since(startTime).Milliseconds()
	if execErr != nil {
		log.Status = 2
		log.ErrorMsg = execErr.Error()
		// 失败时自动通知超管
		notifyAdminOnFailure(jobId, jobName, execErr.Error())
	} else {
		log.Status = 1
	}
	log.Create()

	// 更新任务的 lastRunAt 和 nextRunAt
	if entryId, ok := jobEntryMap.Load(jobId); ok {
		entry := scheduler.Entry(entryId.(cron.EntryID))
		nextRun := entry.Next
		now := time.Now()
		(&models.JobTask{}).UpdateRunTime(jobId, &now, &nextRun)
	}
}

// RunJobNow 立即执行一次任务
func RunJobNow(job models.JobTask) {
	go executeJob(job.ID, job.Name, job.Content, job.JobType)
}

// GetNextRunTime 获取任务下次执行时间预览
func GetNextRunTime(cronExpr string) (string, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return "", err
	}
	next := schedule.Next(time.Now())
	return next.Format("2006-01-02 15:04:05"), nil
}

// StopScheduler 停止调度器
func StopScheduler() {
	if scheduler != nil {
		scheduler.Stop()
	}
}

// performHealthCheck 执行所有导航链接的健康检查
// 使用 HEAD 请求，超时 5 秒，连续失败阈值 3 次才报警
func performHealthCheck() error {
	// 获取所有有 URL 的图标（排除维护模式）
	var icons []models.ItemIcon
	if err := global.Db.Where("url != '' AND url IS NOT NULL AND health_muted = 0").
		Select("id, user_id, title, url, lan_url, health_status, health_fail_count").
		Find(&icons).Error; err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: 10,
		},
	}

	const HEALTH_FAIL_THRESHOLD = 3 // 连续失败 3 次才触发报警

	for _, icon := range icons {
		checkUrl := icon.Url
		if icon.LanUrl != "" {
			checkUrl = icon.LanUrl
		}

		if strings.HasPrefix(checkUrl, "192.168.") ||
			strings.HasPrefix(checkUrl, "10.") ||
			strings.HasPrefix(checkUrl, "172.") {
			continue
		}

		resp, err := client.Head(checkUrl)
		status := 1            // 正常
		failCount := 0         // 重置失败计数
		if err != nil || (resp != nil && resp.StatusCode >= 500) {
			failCount = icon.HealthFailCount + 1
			if failCount >= HEALTH_FAIL_THRESHOLD {
				status = 2 // 连续失败达到阈值
			}
		}
		if resp != nil {
			resp.Body.Close()
		}

		global.Db.Model(&models.ItemIcon{}).Where("id=?", icon.ID).Updates(map[string]interface{}{
			"health_status":     status,
			"health_check_at":   now,
			"health_fail_count": failCount,
		})

		// 仅在首次达到阈值时通知（避免重复轰炸）
		if status == 2 && icon.HealthStatus != 2 {
			notifyIconOwner(icon.ID, icon.Title, checkUrl, icon.UserId)
		}
	}

	return nil
}

// notifyIconOwner 通知图标拥有者其导航链接异常
func notifyIconOwner(iconId uint, title string, url string, userId uint) {
	content := fmt.Sprintf("导航链接「%s」健康检查异常，请确认服务是否正常运行。\nURL: %s", title, url)
	notice := models.Notice{
		Title:       fmt.Sprintf("[链接异常] %s", title),
		Content:     content,
		DisplayType: models.NOTICE_DISPLAY_TYPE_HOME,
		NoticeType:  models.NOTICE_TYPE_MESSAGE,
		OneRead:     0,
		Status:      models.NOTICE_STATUS_ENABLED,
		UserId:      0,
		TargetUserIds: fmt.Sprintf(",%d,", userId),
	}
	global.Db.Create(&notice)
}

// notifyAdminOnFailure 任务失败时自动给超级管理员发送站内信
func notifyAdminOnFailure(jobId uint, jobName, errMsg string) {
	// 查找所有超级管理员
	var admins []models.User
	if err := global.Db.Where("role=?", ROLE_SUPER_ADMIN).Find(&admins).Error; err != nil {
		return
	}

	// 截断错误信息，避免过长
	if len(errMsg) > 500 {
		errMsg = errMsg[:500] + "..."
	}

	content := fmt.Sprintf("定时任务「%s」(ID:%d) 执行失败：\n%s", jobName, jobId, errMsg)

	for _, admin := range admins {
		notice := models.Notice{
			Title:       fmt.Sprintf("[任务异常] %s", jobName),
			Content:     content,
			DisplayType: models.NOTICE_DISPLAY_TYPE_HOME,
			NoticeType:  models.NOTICE_TYPE_MESSAGE,
			OneRead:     0, // 站内信，不一次性阅读
			Status:      models.NOTICE_STATUS_ENABLED,
			UserId:      0, // 系统触发
			TargetUserIds: fmt.Sprintf(",%d,", admin.ID),
		}
		if err := global.Db.Create(&notice).Error; err != nil {
			global.Logger.Debug("任务失败通知发送失败:", err.Error())
		}
	}
}

const (
	ROLE_SUPER_ADMIN = 1
	JOB_TYPE_REMINDER   = 1
	JOB_TYPE_HEALTH_CHECK = 3 // 导航链接健康检查
)

// 使用空 context 避免 import 问题
var _ = context.Background