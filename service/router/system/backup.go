package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitBackupRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.BackupApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/backup/create", api.CreateBackup)
		r.POST("/backup/list", api.GetList)
		r.POST("/backup/get", api.GetById)
		r.POST("/backup/restore", api.Restore)
		r.POST("/backup/delete", api.Delete)
		r.POST("/backup/export", api.Export)
		r.POST("/backup/import", api.Import)
		r.POST("/backup/task/list", api.GetTaskList)
		r.POST("/backup/task/create", api.CreateTask)
		r.POST("/backup/task/update", api.UpdateTask)
		r.POST("/backup/task/delete", api.DeleteTask)
	}
}