package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitPermissionRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.PermissionApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/permission/list", middleware.PermissionInterceptor("permission:view"), api.GetList)
		r.POST("/permission/byModule", middleware.PermissionInterceptor("permission:view"), api.GetByModule)
		r.POST("/permission/matrix", middleware.PermissionInterceptor("permission:view"), api.GetPermissionMatrix)
		r.POST("/permission/saveRolePermissions", middleware.PermissionInterceptor("permission:edit"), api.SaveRolePermissions)
		r.POST("/permission/getRolePermissions", middleware.PermissionInterceptor("permission:view"), api.GetRolePermissions)
		r.POST("/permission/create", middleware.PermissionInterceptor("permission:edit"), api.Create)
		r.POST("/permission/update", middleware.PermissionInterceptor("permission:edit"), api.Update)
	}
}