package router

import (
	"os"
	"path/filepath"
	"strings"

	"sun-panel/api/api_v1/middleware"
	"sun-panel/global"
	// "sun-panel/router/admin"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// gzip 压缩（全局）
	routerGroup.Use(middleware.GzipInterceptor())

	// 操作日志中间件（全局，异步不阻塞）
	routerGroup.Use(middleware.OperationLogInterceptor)

	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)
	openness.Init(routerGroup)

	// WEB文件服务
	{
		webPath := "./web"
		router.Use(middleware.CacheControlInterceptor())
		router.Use(middleware.GzipInterceptor())
		router.StaticFile("/", webPath+"/index.html")
		router.Static("/assets", webPath+"/assets")
		router.Static("/custom", webPath+"/custom")
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")

		// PWA 静态产物：手机「添加到主屏幕」与离线缓存所需
		router.GET("/manifest.webmanifest", func(c *gin.Context) {
			c.Header("Content-Type", "application/manifest+json")
			c.File(webPath + "/manifest.webmanifest")
		})
		router.StaticFile("/sw.js", webPath+"/sw.js")
		router.StaticFile("/registerSW.js", webPath+"/registerSW.js")
		router.StaticFile("/pwa-192x192.png", webPath+"/pwa-192x192.png")
		router.StaticFile("/pwa-512x512.png", webPath+"/pwa-512x512.png")

		// SPA fallback：真实文件优先服务（覆盖 workbox-*.js 等构建产物），否则返回 index.html
		router.NoRoute(func(c *gin.Context) {
			rel := c.Request.URL.Path
			abs := filepath.Join(webPath, filepath.Clean("/"+rel))
			if strings.HasPrefix(abs, filepath.Clean(webPath)) {
				if info, err := os.Stat(abs); err == nil && !info.IsDir() {
					c.File(abs)
					return
				}
			}
			c.File(webPath + "/index.html")
		})
	}

	// 上传的文件
	sourcePath := global.Config.GetValueString("base", "source_path")
	router.Static(sourcePath[1:], sourcePath)

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
