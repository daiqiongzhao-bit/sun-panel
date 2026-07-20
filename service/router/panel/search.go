package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitSearchRouter(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.SearchApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/search/local", api.LocalSearch)
	}
}