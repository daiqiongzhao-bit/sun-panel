package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitNoticeManageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.NoticeManageApi

	// 管理员接口
	rAdmin := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		rAdmin.POST("/noticeManage/getList", api.GetList)
		rAdmin.POST("/noticeManage/get", api.GetById)
		rAdmin.POST("/noticeManage/create", api.Create)
		rAdmin.POST("/noticeManage/update", api.Update)
		rAdmin.POST("/noticeManage/delete", api.Delete)
	}

	// 登录用户可见通知（公告弹窗用）
	rUser := router.Group("", middleware.LoginInterceptor)
	{
		rUser.POST("/noticeManage/getVisibleNotices", api.GetVisibleNotices)
	}
}