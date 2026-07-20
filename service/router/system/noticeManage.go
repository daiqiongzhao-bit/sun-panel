package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitNoticeManageRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.NoticeManageApi

	// 公告管理（按细粒度权限校验）
	rAdmin := router.Group("", middleware.LoginInterceptor)
	{
		rAdmin.POST("/noticeManage/getList", middleware.PermissionInterceptor("notice:view"), api.GetList)
		rAdmin.POST("/noticeManage/get", middleware.PermissionInterceptor("notice:view"), api.GetById)
		rAdmin.POST("/noticeManage/create", middleware.PermissionInterceptor("notice:create"), api.Create)
		rAdmin.POST("/noticeManage/update", middleware.PermissionInterceptor("notice:edit"), api.Update)
		rAdmin.POST("/noticeManage/delete", middleware.PermissionInterceptor("notice:delete"), api.Delete)
	}

	// 登录用户可见通知（公告弹窗用，所有登录用户可读）
	rUser := router.Group("", middleware.LoginInterceptor)
	{
		rUser.POST("/noticeManage/getVisibleNotices", api.GetVisibleNotices)
	}
}