package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitDepartmentRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.DepartmentApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/department/list", api.GetList)
		r.POST("/department/tree", api.GetTreeList)
		r.POST("/department/get", api.GetById)
		r.POST("/department/create", api.Create)
		r.POST("/department/update", api.Update)
		r.POST("/department/delete", api.Delete)
	}
}