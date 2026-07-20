package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitBackupRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.BackupApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/backup/create", middleware.PermissionInterceptor("backup:create"), api.CreateBackup)
		r.POST("/backup/list", middleware.PermissionInterceptor("backup:view"), api.GetList)
		r.POST("/backup/get", middleware.PermissionInterceptor("backup:view"), api.GetById)
		r.POST("/backup/restore", middleware.PermissionInterceptor("backup:restore"), api.Restore)
		r.POST("/backup/delete", middleware.PermissionInterceptor("backup:delete"), api.Delete)
		r.POST("/backup/export", middleware.PermissionInterceptor("backup:create"), api.Export)
		r.POST("/backup/import", middleware.PermissionInterceptor("backup:restore"), api.Import)
		r.POST("/backup/task/list", middleware.PermissionInterceptor("backup:view"), api.GetTaskList)
		r.POST("/backup/task/create", middleware.PermissionInterceptor("backup:create"), api.CreateTask)
		r.POST("/backup/task/update", middleware.PermissionInterceptor("backup:create"), api.UpdateTask)
		r.POST("/backup/task/delete", middleware.PermissionInterceptor("backup:delete"), api.DeleteTask)
	}
}