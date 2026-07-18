package system

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BackupApi struct{}

const (
	BACKUP_MODE_DB    = 1
	BACKUP_MODE_FULL  = 2
)

func (a *BackupApi) CreateBackup(c *gin.Context) {
	type Request struct {
		Mode int    `json:"mode"`
		Name string `json:"name"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Mode != BACKUP_MODE_DB && req.Mode != BACKUP_MODE_FULL {
		apiReturn.ErrorParamFomat(c, "备份模式只能是1(数据库备份)或2(全部数据备份)")
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	backupDir := global.Config.GetValueString("base", "backup_path")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.MkdirAll(backupDir, os.ModePerm)
	}

	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s", timestamp, req.Name)
	if req.Name == "" {
		fileName = timestamp
	}

	var filePath string
	var backupSize int64
	var err error

	// 获取源数据库路径
	sourceDbPath := global.Config.GetValueString("sqlite", "file_path")
	if sourceDbPath == "" {
		sourceDbPath = "database.db"
	}

	if req.Mode == BACKUP_MODE_DB {
		filePath = fmt.Sprintf("%s/%s.sql", backupDir, fileName)
		err = models.BackupDatabase(sourceDbPath, filePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "备份数据库失败: "+err.Error())
			return
		}
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "获取备份文件信息失败: "+err.Error())
			return
		}
		backupSize = fileInfo.Size()
	} else {
		filePath = fmt.Sprintf("%s/%s.zip", backupDir, fileName)
		err = models.BackupFullData(sourceDbPath, filePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "备份全部数据失败: "+err.Error())
			return
		}
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "获取备份文件信息失败: "+err.Error())
			return
		}
		backupSize = fileInfo.Size()
	}

	backup := models.Backup{
		Name:        req.Name,
		Mode:        req.Mode,
		FilePath:    filePath,
		Size:        backupSize,
		CreatorId:   userInfo.ID,
		Status:      1,
	}

	created, err := backup.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *BackupApi) GetList(c *gin.Context) {
	type Request struct {
		PageNum  int `json:"pageNum"`
		PageSize int `json:"pageSize"`
		Mode     int `json:"mode"`
	}
	req := Request{PageNum: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mBackup := models.Backup{}
	list, total, err := mBackup.GetList(req.PageNum, req.PageSize, req.Mode)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (a *BackupApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mBackup := models.Backup{}
	backup, err := mBackup.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, backup)
}

func (a *BackupApi) Restore(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mBackup := models.Backup{}
	backup, err := mBackup.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	if backup.Mode == BACKUP_MODE_DB {
		err = models.RestoreDatabase(backup.FilePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "恢复数据库失败: "+err.Error())
			return
		}
	} else {
		err = models.RestoreFullData(backup.FilePath)
		if err != nil {
			apiReturn.ErrorDatabase(c, "恢复全部数据失败: "+err.Error())
			return
		}
	}

	err = mBackup.UpdateStatus(req.Id, 2)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *BackupApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mBackup := models.Backup{}
	backup, err := mBackup.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	os.Remove(backup.FilePath)

	err = mBackup.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *BackupApi) Export(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mBackup := models.Backup{}
	backup, err := mBackup.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	fileName := path.Base(backup.FilePath)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(backup.FilePath)
}

func (a *BackupApi) Import(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	backupDir := global.Config.GetValueString("base", "backup_path")
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := path.Ext(f.Filename)
	if fileExt != ".sql" && fileExt != ".zip" {
		apiReturn.ErrorParamFomat(c, "只支持.sql和.zip格式的备份文件")
		return
	}

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.MkdirAll(backupDir, os.ModePerm)
	}

	filePath := fmt.Sprintf("%s/%s", backupDir, filepath.Base(f.Filename))
	c.SaveUploadedFile(f, filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		apiReturn.ErrorParamFomat(c, "获取导入文件信息失败: "+err.Error())
		return
	}

	backup := models.Backup{
		Name:      f.Filename,
		Mode:      BACKUP_MODE_DB,
		FilePath:  filePath,
		Size:      fileInfo.Size(),
		CreatorId: userInfo.ID,
		Status:    1,
	}
	if fileExt == ".zip" {
		backup.Mode = BACKUP_MODE_FULL
	}

	_, err = backup.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *BackupApi) GetTaskList(c *gin.Context) {
	mTask := models.BackupTask{}
	list, err := mTask.GetList()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *BackupApi) CreateTask(c *gin.Context) {
	type Request struct {
		Name         string `json:"name" validate:"required"`
		Mode         int    `json:"mode"`
		CronExpr     string `json:"cronExpr" validate:"required"`
		RetentionDays int   `json:"retentionDays"`
		Status       int    `json:"status"`
	}
	req := Request{Mode: BACKUP_MODE_DB, RetentionDays: 7, Status: 1}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	if req.Mode != BACKUP_MODE_DB && req.Mode != BACKUP_MODE_FULL {
		apiReturn.ErrorParamFomat(c, "备份模式只能是1(数据库备份)或2(全部数据备份)")
		return
	}

	task := models.BackupTask{
		Name:          req.Name,
		Mode:          req.Mode,
		CronExpr:      req.CronExpr,
		RetentionDays: req.RetentionDays,
		Status:        req.Status,
	}

	created, err := task.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	if req.Status == 1 {
		err = models.RegisterCronTask(created.ID, created.CronExpr)
		if err != nil {
			apiReturn.ErrorByCodeAndMsg(c, 1003, "注册定时任务失败: "+err.Error())
			return
		}
	}

	apiReturn.SuccessData(c, created)
}

func (a *BackupApi) UpdateTask(c *gin.Context) {
	type Request struct {
		Id            uint   `json:"id"`
		Name          string `json:"name" validate:"required"`
		Mode          int    `json:"mode"`
		CronExpr      string `json:"cronExpr" validate:"required"`
		RetentionDays int    `json:"retentionDays"`
		Status        int    `json:"status"`
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

	task := models.BackupTask{
		Name:          req.Name,
		Mode:          req.Mode,
		CronExpr:      req.CronExpr,
		RetentionDays: req.RetentionDays,
		Status:        req.Status,
	}

	err := task.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	if req.Status == 1 {
		err = models.RegisterCronTask(req.Id, req.CronExpr)
		if err != nil {
			apiReturn.ErrorByCodeAndMsg(c, 1003, "注册定时任务失败: "+err.Error())
			return
		}
	} else {
		err = models.UnregisterCronTask(req.Id)
		if err != nil {
			apiReturn.ErrorByCodeAndMsg(c, 1003, "注销定时任务失败: "+err.Error())
			return
		}
	}

	apiReturn.Success(c)
}

func (a *BackupApi) DeleteTask(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	err := models.UnregisterCronTask(req.Id)
	if err != nil {
		apiReturn.ErrorByCodeAndMsg(c, 1003, "注销定时任务失败: "+err.Error())
		return
	}

	mTask := models.BackupTask{}
	err = mTask.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}