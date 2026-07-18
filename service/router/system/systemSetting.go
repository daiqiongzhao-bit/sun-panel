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

	// 管理端读写（需要管理员权限）
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/systemSetting/logo/set", api.SetLogoConfig)
		r.POST("/systemSetting/logo/upload", api.UploadLogo)
		r.POST("/systemSetting/background/set", api.SetBackgroundConfig)
		r.POST("/systemSetting/background/upload", api.UploadBackground)
	}
}