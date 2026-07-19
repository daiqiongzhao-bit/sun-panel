package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/lib/job"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type JobApi struct{}

// GetList 获取任务列表
func (a *JobApi) GetList(c *gin.Context) {
	type Request struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Keyword  string `json:"keyword"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mJob := models.JobTask{}
	list, count, err := mJob.GetList(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

// Create 创建任务
func (a *JobApi) Create(c *gin.Context) {
	type Request struct {
		Name     string `json:"name" validate:"required,min=1,max=100"`
		JobType  int    `json:"jobType"`
		CronExpr string `json:"cronExpr" validate:"required"`
		Content  string `json:"content" validate:"required"`
	}
	req := Request{JobType: models.JOB_TYPE_REMINDER}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// 验证 Cron 表达式
	nextTime, err := job.GetNextRunTime(req.CronExpr)
	if err != nil {
		apiReturn.ErrorParamFomat(c, "Cron 表达式无效: "+err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	task := models.JobTask{
		Name:      req.Name,
		JobType:   req.JobType,
		CronExpr:  req.CronExpr,
		Content:   req.Content,
		Status:    models.JOB_STATUS_RUNNING,
		CreatorId: userInfo.ID,
	}

	created, err := task.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 注册到调度器
	job.RegisterJob(created)

	apiReturn.SuccessData(c, gin.H{"id": created.ID, "nextRunAt": nextTime})
}

// Update 更新任务
func (a *JobApi) Update(c *gin.Context) {
	type Request struct {
		Id       uint   `json:"id" validate:"required"`
		Name     string `json:"name" validate:"required,min=1,max=100"`
		JobType  int    `json:"jobType"`
		CronExpr string `json:"cronExpr" validate:"required"`
		Content  string `json:"content" validate:"required"`
		Status   int    `json:"status"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	// 验证 Cron 表达式
	_, err := job.GetNextRunTime(req.CronExpr)
	if err != nil {
		apiReturn.ErrorParamFomat(c, "Cron 表达式无效: "+err.Error())
		return
	}

	task := models.JobTask{
		Name:     req.Name,
		JobType:  req.JobType,
		CronExpr: req.CronExpr,
		Content:  req.Content,
		Status:   req.Status,
	}

	if err := task.Update(req.Id); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 重新注册到调度器
	job.UnregisterJob(req.Id)
	if req.Status == models.JOB_STATUS_RUNNING {
		updated, _ := (&models.JobTask{}).GetById(req.Id)
		job.RegisterJob(updated)
	}

	apiReturn.Success(c)
}

// Pause 暂停任务
func (a *JobApi) Pause(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	job.UnregisterJob(req.Id)
	(&models.JobTask{}).UpdateStatus(req.Id, models.JOB_STATUS_PAUSED)

	apiReturn.Success(c)
}

// Start 启动任务
func (a *JobApi) Start(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	task, err := (&models.JobTask{}).GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, "任务不存在")
		return
	}

	task.Status = models.JOB_STATUS_RUNNING
	task.Update(req.Id)
	job.RegisterJob(task)

	apiReturn.Success(c)
}

// RunNow 立即执行一次
func (a *JobApi) RunNow(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	task, err := (&models.JobTask{}).GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, "任务不存在")
		return
	}

	job.RunJobNow(task)
	apiReturn.Success(c)
}

// Delete 删除任务
func (a *JobApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	job.UnregisterJob(req.Id)
	(&models.JobTask{}).Delete(req.Id)
	(&models.JobLog{}).DeleteByJobId(req.Id)

	apiReturn.Success(c)
}

// GetLogList 获取任务执行日志
func (a *JobApi) GetLogList(c *gin.Context) {
	type Request struct {
		JobId    uint `json:"jobId"`
		Page     int  `json:"page"`
		PageSize int  `json:"pageSize"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mLog := models.JobLog{}
	list, count, err := mLog.GetList(req.JobId, req.Page, req.PageSize)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

// PreviewCron 预览 Cron 表达式的下次执行时间
func (a *JobApi) PreviewCron(c *gin.Context) {
	type Request struct {
		CronExpr string `json:"cronExpr"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	nextTime, err := job.GetNextRunTime(req.CronExpr)
	if err != nil {
		apiReturn.ErrorParamFomat(c, "Cron 表达式无效: "+err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"nextRunAt": nextTime})
}