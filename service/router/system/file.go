package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitFileRouter(router *gin.RouterGroup) {
	FileApi := api_v1.ApiGroupApp.ApiSystem.FileApi

	// 验证项目的权限(有访问密码的需要验证访问token)
	private := router.Group("", middleware.LoginInterceptor)
	{
		private.POST("/file/uploadImg", middleware.PermissionInterceptor("file:upload"), FileApi.UploadImg)
		private.POST("/file/uploadFiles", middleware.PermissionInterceptor("file:upload"), FileApi.UploadFiles)

		private.POST("/file/getList", middleware.PermissionInterceptor("file:view"), FileApi.GetList)
		private.POST("/file/deletes", middleware.PermissionInterceptor("file:delete"), FileApi.Deletes)

	}

}
