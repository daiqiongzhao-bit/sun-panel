package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitGalleryRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.GalleryApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/gallery/list", api.GetList)
		r.POST("/gallery/get", api.GetById)
		r.POST("/gallery/upload", api.Upload)
		r.POST("/gallery/update", api.Update)
		r.POST("/gallery/delete", api.Delete)
		r.POST("/gallery/batchDelete", api.BatchDelete)
		r.POST("/gallery/category/list", api.GetCategories)
		r.POST("/gallery/category/create", api.CreateCategory)
		r.POST("/gallery/category/update", api.UpdateCategory)
		r.POST("/gallery/category/delete", api.DeleteCategory)
	}
}