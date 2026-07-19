package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitLogRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.LogApi

	rAdmin := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		rAdmin.POST("/log/operationList", api.GetOperationLogList)
		rAdmin.POST("/log/loginList", api.GetLoginLogList)
		rAdmin.POST("/log/clearOperation", api.ClearOperationLog)
		rAdmin.POST("/log/clearLogin", api.ClearLoginLog)
	}
}