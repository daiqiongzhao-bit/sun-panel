package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitLogin(router *gin.RouterGroup) {
	loginApi := api_v1.ApiGroupApp.ApiSystem.LoginApi

	router.POST("/login", loginApi.Login)
	router.POST("/login/2fa", loginApi.Login2FA)
	router.POST("/logout", middleware.LoginInterceptor, loginApi.Logout)

	// 两步验证(2FA)自助管理，需登录
	router.POST("/twofa/status", middleware.LoginInterceptor, loginApi.TwoFAStatus)
	router.POST("/twofa/enable", middleware.LoginInterceptor, loginApi.TwoFAEnable)
	router.POST("/twofa/confirm", middleware.LoginInterceptor, loginApi.TwoFAConfirm)
	router.POST("/twofa/disable", middleware.LoginInterceptor, loginApi.TwoFADisable)

}
