package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitIconRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiSystem.IconApi
	r := router.Group("", middleware.LoginInterceptor, middleware.AdminInterceptor)
	{
		r.POST("/icon/get", api.GetById)
		r.POST("/icon/getByName", api.GetByName)
		r.POST("/icon/list", api.GetList)
		r.POST("/icon/batch", api.GetBatch)
		r.POST("/icon/create", api.Create)
		r.POST("/icon/update", api.Update)
		r.POST("/icon/delete", api.Delete)
		r.POST("/icon/favorite", api.Favorite)
		r.POST("/icon/favorites", api.GetFavorites)
		r.POST("/icon/category/list", api.GetCategories)
		r.POST("/icon/category/create", api.CreateCategory)
		r.POST("/icon/category/update", api.UpdateCategory)
		r.POST("/icon/category/delete", api.DeleteCategory)
	}
}