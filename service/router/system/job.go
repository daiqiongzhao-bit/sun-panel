package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitJobRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.JobApi

	rAdmin := router.Group("", middleware.LoginInterceptor)
	{
		rAdmin.POST("/job/getList", middleware.PermissionInterceptor("task:view"), api.GetList)
		rAdmin.POST("/job/create", middleware.PermissionInterceptor("task:create"), api.Create)
		rAdmin.POST("/job/update", middleware.PermissionInterceptor("task:edit"), api.Update)
		rAdmin.POST("/job/pause", middleware.PermissionInterceptor("task:edit"), api.Pause)
		rAdmin.POST("/job/start", middleware.PermissionInterceptor("task:edit"), api.Start)
		rAdmin.POST("/job/runNow", middleware.PermissionInterceptor("task:edit"), api.RunNow)
		rAdmin.POST("/job/delete", middleware.PermissionInterceptor("task:delete"), api.Delete)
		rAdmin.POST("/job/getLogList", middleware.PermissionInterceptor("task:view"), api.GetLogList)
		rAdmin.POST("/job/previewCron", middleware.PermissionInterceptor("task:view"), api.PreviewCron)
	}
}