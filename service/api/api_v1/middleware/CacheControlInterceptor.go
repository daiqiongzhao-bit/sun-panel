package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CacheControlInterceptor 为带内容哈希的静态资源(/assets、/custom)设置长期不可变缓存，
// 文件名含内容 hash，内容变化即换名，因此可安全设为 1 年 immutable，使重复访问秒开(零网络往返)。
// 对 /api 等动态接口不设置缓存头，避免误缓存。
func CacheControlInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/assets/") || strings.HasPrefix(p, "/custom/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.Next()
	}
}
