package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitImportExportRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.ImportExportApi
	// 需要管理员权限
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/system/exportUsers", api.ExportUsers)
		r.POST("/system/importUsers", api.ImportUsers)
		r.POST("/system/exportRoles", api.ExportRoles)
		r.POST("/system/importRoles", api.ImportRoles)
	}
}