package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitIconRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.IconApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/icon/get", middleware.PermissionInterceptor("icon:view"), api.GetById)
		r.POST("/icon/getByName", middleware.PermissionInterceptor("icon:view"), api.GetByName)
		r.POST("/icon/list", middleware.PermissionInterceptor("icon:view"), api.GetList)
		r.POST("/icon/batch", middleware.PermissionInterceptor("icon:view"), api.GetBatch)
		r.POST("/icon/create", middleware.PermissionInterceptor("icon:upload"), api.Create)
		r.POST("/icon/update", middleware.PermissionInterceptor("icon:upload"), api.Update)
		r.POST("/icon/delete", middleware.PermissionInterceptor("icon:delete"), api.Delete)
		r.POST("/icon/favorite", middleware.PermissionInterceptor("icon:view"), api.Favorite)
		r.POST("/icon/favorites", middleware.PermissionInterceptor("icon:view"), api.GetFavorites)
		r.POST("/icon/category/list", middleware.PermissionInterceptor("icon:view"), api.GetCategories)
		r.POST("/icon/category/create", middleware.PermissionInterceptor("icon:upload"), api.CreateCategory)
		r.POST("/icon/category/update", middleware.PermissionInterceptor("icon:upload"), api.UpdateCategory)
		r.POST("/icon/category/delete", middleware.PermissionInterceptor("icon:delete"), api.DeleteCategory)
	}
}