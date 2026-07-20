package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitLogRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.LogApi

	rAdmin := router.Group("", middleware.LoginInterceptor)
	{
		rAdmin.POST("/log/operationList", middleware.PermissionInterceptor("operation_log:view"), api.GetOperationLogList)
		rAdmin.POST("/log/loginList", middleware.PermissionInterceptor("login_log:view"), api.GetLoginLogList)
		rAdmin.POST("/log/clearOperation", middleware.PermissionInterceptor("operation_log:clear"), api.ClearOperationLog)
		rAdmin.POST("/log/clearLogin", middleware.PermissionInterceptor("login_log:clear"), api.ClearLoginLog)
	}
}