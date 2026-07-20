package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitStickyNoteRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.StickyNoteApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/stickyNote/getList", api.GetList)
		r.POST("/panel/stickyNote/create", api.Create)
		r.POST("/panel/stickyNote/update", api.Update)
		r.POST("/panel/stickyNote/delete", api.Delete)
	}
}

func InitPasteBinRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.PasteBinApi

	// 需要登录
	rUser := router.Group("", middleware.LoginInterceptor)
	{
		rUser.POST("/panel/pasteBin/create", api.Create)
		rUser.POST("/panel/pasteBin/update", api.Update)
		rUser.POST("/panel/pasteBin/myList", api.GetMyList)
		rUser.POST("/panel/pasteBin/delete", api.Delete)
		rUser.POST("/panel/pasteBin/accessUrl", api.GetAccessUrl)
		rUser.POST("/panel/pasteBin/uploadFile", api.UploadFile)
	}

	// 无需登录（通过 code 访问）
	router.POST("/panel/pasteBin/getByCode", api.GetByCode)
	router.POST("/panel/pasteBin/downloadFile", api.DownloadFile)
	router.GET("/panel/pasteBin/downloadFile", api.DownloadFile)
}