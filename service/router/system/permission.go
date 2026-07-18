package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitPermissionRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.PermissionApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/permission/list", api.GetList)
		r.POST("/permission/byModule", api.GetByModule)
		r.POST("/permission/matrix", api.GetPermissionMatrix)
		r.POST("/permission/saveRolePermissions", api.SaveRolePermissions)
		r.POST("/permission/getRolePermissions", api.GetRolePermissions)
		r.POST("/permission/create", api.Create)
		r.POST("/permission/update", api.Update)
	}
}