package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitSystemSettingRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.SystemSettingApi
	
	// 公开读取（无需登录，供登录页/首页展示）
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/systemSetting/logo/get", api.GetLogoConfig)
		rPublic.POST("/systemSetting/background/get", api.GetBackgroundConfig)
		rPublic.POST("/systemSetting/background/preset", api.GetPresetBackgrounds)
	}

	// 管理端读写（按细粒度权限校验）
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/systemSetting/logo/set", middleware.PermissionInterceptor("setting:edit"), api.SetLogoConfig)
		r.POST("/systemSetting/logo/upload", middleware.PermissionInterceptor("setting:edit"), api.UploadLogo)
		r.POST("/systemSetting/background/set", middleware.PermissionInterceptor("setting:edit"), api.SetBackgroundConfig)
		r.POST("/systemSetting/background/upload", middleware.PermissionInterceptor("setting:edit"), api.UploadBackground)
	}
}