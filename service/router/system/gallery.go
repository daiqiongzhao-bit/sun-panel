package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitGalleryRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.GalleryApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/gallery/list", middleware.PermissionInterceptor("gallery:view"), api.GetList)
		r.POST("/gallery/get", middleware.PermissionInterceptor("gallery:view"), api.GetById)
		r.POST("/gallery/upload", middleware.PermissionInterceptor("gallery:upload"), api.Upload)
		r.POST("/gallery/update", middleware.PermissionInterceptor("gallery:upload"), api.Update)
		r.POST("/gallery/delete", middleware.PermissionInterceptor("gallery:delete"), api.Delete)
		r.POST("/gallery/batchDelete", middleware.PermissionInterceptor("gallery:delete"), api.BatchDelete)
		r.POST("/gallery/category/list", middleware.PermissionInterceptor("gallery:view"), api.GetCategories)
		r.POST("/gallery/category/create", middleware.PermissionInterceptor("gallery:upload"), api.CreateCategory)
		r.POST("/gallery/category/update", middleware.PermissionInterceptor("gallery:upload"), api.UpdateCategory)
		r.POST("/gallery/category/delete", middleware.PermissionInterceptor("gallery:delete"), api.DeleteCategory)
	}
}