package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitItemIcon(router *gin.RouterGroup) {
	itemIcon := api_v1.ApiGroupApp.ApiPanel.ItemIcon
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/itemIcon/edit", middleware.PermissionInterceptorAny("icon:edit", "icon:create", "icon:upload"), itemIcon.Edit)
		r.POST("/panel/itemIcon/deletes", middleware.PermissionInterceptor("icon:delete"), itemIcon.Deletes)
		r.POST("/panel/itemIcon/saveSort", middleware.PermissionInterceptor("icon:edit"), itemIcon.SaveSort)
		r.POST("/panel/itemIcon/addMultiple", middleware.PermissionInterceptorAny("icon:create", "icon:upload"), itemIcon.AddMultiple)
		r.POST("/panel/itemIcon/getSiteFavicon", middleware.PermissionInterceptor("icon:view"), itemIcon.GetSiteFavicon)
	}

	// 公开模式
	rPublic := router.Group("", middleware.PublicModeInterceptor)
	{
		rPublic.POST("/panel/itemIcon/getListByGroupId", itemIcon.GetListByGroupId)
	}
}
