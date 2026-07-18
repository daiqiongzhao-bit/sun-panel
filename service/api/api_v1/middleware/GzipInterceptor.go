package middleware

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (w gzipResponseWriter) WriteString(s string) (int, error) {
	return w.Writer.Write([]byte(s))
}

func (w gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func GzipInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		writer := gzip.NewWriter(c.Writer)

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		gzWriter := gzipResponseWriter{Writer: writer, ResponseWriter: c.Writer}
		c.Writer = gzWriter
		c.Next()

		writer.Close()
	}
}
