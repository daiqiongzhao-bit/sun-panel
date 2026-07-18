package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitJobRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.JobApi

	rAdmin := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		rAdmin.POST("/job/getList", api.GetList)
		rAdmin.POST("/job/create", api.Create)
		rAdmin.POST("/job/update", api.Update)
		rAdmin.POST("/job/pause", api.Pause)
		rAdmin.POST("/job/start", api.Start)
		rAdmin.POST("/job/runNow", api.RunNow)
		rAdmin.POST("/job/delete", api.Delete)
		rAdmin.POST("/job/getLogList", api.GetLogList)
		rAdmin.POST("/job/previewCron", api.PreviewCron)
	}
}