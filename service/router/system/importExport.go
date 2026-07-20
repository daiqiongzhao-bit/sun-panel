package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitImportExportRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.ImportExportApi
	// 按细粒度权限校验（导入导出模块）
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/system/exportUsers", middleware.PermissionInterceptor("import_export:export"), api.ExportUsers)
		r.POST("/system/importUsers", middleware.PermissionInterceptor("import_export:import"), api.ImportUsers)
		r.POST("/system/exportRoles", middleware.PermissionInterceptor("import_export:export"), api.ExportRoles)
		r.POST("/system/importRoles", middleware.PermissionInterceptor("import_export:import"), api.ImportRoles)
	}
}