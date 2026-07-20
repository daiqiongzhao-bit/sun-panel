package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.RoleApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/role/list", middleware.PermissionInterceptor("role:view"), api.GetList)
		r.POST("/role/get", middleware.PermissionInterceptor("role:view"), api.GetById)
		r.POST("/role/create", middleware.PermissionInterceptor("role:create"), api.Create)
		r.POST("/role/update", middleware.PermissionInterceptor("role:edit"), api.Update)
		r.POST("/role/delete", middleware.PermissionInterceptor("role:delete"), api.Delete)
	}
}