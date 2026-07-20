package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitUsersRouter(router *gin.RouterGroup) {
	userApi := api_v1.ApiGroupApp.ApiPanel.UsersApi

	rAdmin := router.Group("", middleware.LoginInterceptor)
	{
		rAdmin.POST("panel/users/create", middleware.PermissionInterceptor("user:create"), userApi.Create)
		rAdmin.POST("panel/users/update", middleware.PermissionInterceptor("user:edit"), userApi.Update)
		rAdmin.POST("panel/users/getList", middleware.PermissionInterceptor("user:view"), userApi.GetList)
		rAdmin.POST("panel/users/deletes", middleware.PermissionInterceptor("user:delete"), userApi.Deletes)
		rAdmin.POST("panel/users/getPublicVisitUser", middleware.PermissionInterceptor("user:view"), userApi.GetPublicVisitUser)
		rAdmin.POST("panel/users/setPublicVisitUser", middleware.PermissionInterceptor("user:edit"), userApi.SetPublicVisitUser)
	}
}
