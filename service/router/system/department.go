package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitDepartmentRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.DepartmentApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/department/list", middleware.PermissionInterceptor("department:view"), api.GetList)
		r.POST("/department/tree", middleware.PermissionInterceptor("department:view"), api.GetTreeList)
		r.POST("/department/get", middleware.PermissionInterceptor("department:view"), api.GetById)
		r.POST("/department/create", middleware.PermissionInterceptor("department:create"), api.Create)
		r.POST("/department/update", middleware.PermissionInterceptor("department:edit"), api.Update)
		r.POST("/department/delete", middleware.PermissionInterceptor("department:delete"), api.Delete)
	}
}