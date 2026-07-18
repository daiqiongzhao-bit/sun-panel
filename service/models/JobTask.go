package models

import "time"

const (
	JOB_STATUS_RUNNING = 1 // 运行中
	JOB_STATUS_PAUSED  = 2 // 暂停
)

const (
	JOB_TYPE_REMINDER = 1 // 定时提醒（触发时创建通知）
	JOB_TYPE_CUSTOM   = 2 // 自定义（预留扩展）
)

// JobTask 定时任务
type JobTask struct {
	BaseModel
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	JobType     int        `gorm:"type:tinyint(1);default:1" json:"jobType"`
	CronExpr    string     `gorm:"type:varchar(100)" json:"cronExpr"`
	Content     string     `gorm:"type:text" json:"content"`       // 提醒内容/通知内容
	Status      int        `gorm:"type:tinyint(1);default:1" json:"status"`
	LastRunAt   *time.Time `gorm:"index" json:"lastRunAt"`
	NextRunAt   *time.Time `json:"nextRunAt"`
	CreatorId   uint       `gorm:"index" json:"creatorId"`
}

func (m *JobTask) GetList(page, pageSize int, keyword string) ([]JobTask, int64, error) {
	var list []JobTask
	var count int64
	offset := (page - 1) * pageSize

	db := Db.Model(m)
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
	return list, count, err
}

func (m *JobTask) GetAll() ([]JobTask, error) {
	var list []JobTask
	err := Db.Where("status=?", JOB_STATUS_RUNNING).Find(&list).Error
	return list, err
}

func (m *JobTask) GetById(id uint) (JobTask, error) {
	var task JobTask
	err := Db.Where("id=?", id).First(&task).Error
	return task, err
}

func (m *JobTask) Create() (JobTask, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *JobTask) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":       m.Name,
		"job_type":   m.JobType,
		"cron_expr":  m.CronExpr,
		"content":    m.Content,
		"status":     m.Status,
		"last_run_at": m.LastRunAt,
		"next_run_at": m.NextRunAt,
	}).Error
}

func (m *JobTask) UpdateStatus(id uint, status int) error {
	return Db.Model(m).Where("id=?", id).Update("status", status).Error
}

func (m *JobTask) UpdateRunTime(id uint, lastRun, nextRun *time.Time) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"last_run_at": lastRun,
		"next_run_at": nextRun,
	}).Error
}

func (m *JobTask) Delete(id uint) error {
	return Db.Delete(m, id).Error
}

// JobLog 任务执行日志
type JobLog struct {
	BaseModel
	JobId      uint   `gorm:"index" json:"jobId"`
	JobName    string `gorm:"type:varchar(100)" json:"jobName"`
	StartTime  string `gorm:"type:varchar(30)" json:"startTime"`
	Duration   int64  `gorm:"type:bigint" json:"duration"`  // 耗时ms
	Status     int    `gorm:"type:tinyint(1)" json:"status"` // 1成功 2失败
	ErrorMsg   string `gorm:"type:text" json:"errorMsg"`
}

func (m *JobLog) GetList(jobId uint, page, pageSize int) ([]JobLog, int64, error) {
	var list []JobLog
	var count int64
	offset := (page - 1) * pageSize

	db := Db.Model(m)
	if jobId > 0 {
		db = db.Where("job_id=?", jobId)
	}

	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
	return list, count, err
}

func (m *JobLog) Create() error {
	return Db.Create(m).Error
}

func (m *JobLog) DeleteByJobId(jobId uint) error {
	return Db.Where("job_id=?", jobId).Delete(m).Error
}