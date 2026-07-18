package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.RoleApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/role/list", api.GetList)
		r.POST("/role/get", api.GetById)
		r.POST("/role/create", api.Create)
		r.POST("/role/update", api.Update)
		r.POST("/role/delete", api.Delete)
	}
}