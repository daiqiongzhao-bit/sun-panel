package middleware

import (
	"bytes"
	"io"
	"regexp"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// logRecord 独立结构体，避免在协程中引用 gin.Context
type logRecord struct {
	method       string
	path         string
	ip           string
	userAgent    string
	requestBody  string
	responseCode int
	duration     int64
	userId       uint
	username     string
	module       string
	action       string
}

// OperationLogInterceptor 操作日志中间件
// 关键：所有需要的数据在主线程提取，协程内只使用 logRecord，绝不引用 gin.Context
func OperationLogInterceptor(c *gin.Context) {
	// 公开模式请求不记录操作日志，避免产生大量无意义的 userId=0 记录
	visitMode := base.GetCurrentVisitMode(c)
	if visitMode == base.VISIT_MODE_PUBLIC {
		c.Next()
		return
	}

	startTime := time.Now()

	// 读取请求体（读完后需要重新放回，否则后续绑定会失败）
	var requestBody string
	if c.Request.Body != nil && c.Request.Method != "GET" {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		requestBody = desensitize(string(bodyBytes))
		if len(requestBody) > 2000 {
			requestBody = requestBody[:2000] + "..."
		}
	}

	// ====== 主线程：提前提取不需要用户信息的数据 ======
	rec := &logRecord{
		method:      c.Request.Method,
		path:        c.Request.URL.Path,
		ip:          c.ClientIP(),
		userAgent:   c.Request.UserAgent(),
		requestBody: requestBody,
	}

	// 提取模块名
	pathParts := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")
	if len(pathParts) >= 1 {
		rec.module = pathParts[0]
	}

	// 提取操作类型
	if strings.Contains(rec.path, "/create") || strings.Contains(rec.path, "/edit") || strings.Contains(rec.path, "/add") {
		if rec.method == "DELETE" || strings.Contains(rec.path, "/deletes") {
			rec.action = "delete"
		} else if rec.path != "" {
			rec.action = "create/update"
		}
	} else if strings.Contains(rec.path, "/delete") {
		rec.action = "delete"
	} else if strings.Contains(rec.path, "/login") {
		rec.action = "login"
	} else if strings.Contains(rec.path, "/getList") || strings.Contains(rec.path, "/get") {
		rec.action = "query"
	}

	// 执行请求（LoginInterceptor 在此过程中设置用户信息）
	c.Next()

	// 主线程获取响应状态码和用户信息（LoginInterceptor 之后才能获取）
	rec.responseCode = c.Writer.Status()
	rec.duration = time.Since(startTime).Milliseconds()

	if userInfo, exists := base.GetCurrentUserInfo(c); exists && userInfo.ID != 0 {
		rec.userId = userInfo.ID
		rec.username = userInfo.Username
	}

	// ====== 异步写入：协程内只使用 rec，绝不引用 c ======
	go func(r *logRecord) {
		log := models.OperationLog{
			UserId:       r.userId,
			Username:     r.username,
			Module:       r.module,
			Action:       r.action,
			Method:       r.method,
			Path:         r.path,
			Ip:           r.ip,
			UserAgent:    r.userAgent,
			RequestBody:  r.requestBody,
			ResponseCode: r.responseCode,
			Duration:     r.duration,
		}

		if err := log.Create(); err != nil {
			global.Logger.Debug("操作日志记录失败:", err.Error())
		}
	}(rec)
}

// desensitize 脱敏处理，隐藏 password/token 等敏感字段
func desensitize(body string) string {
	patterns := []string{
		`"password"\s*:\s*"[^"]*"`,
		`"token"\s*:\s*"[^"]*"`,
		`"vcode"\s*:\s*"[^"]*"`,
		`"oldPassword"\s*:\s*"[^"]*"`,
		`"newPassword"\s*:\s*"[^"]*"`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		body = re.ReplaceAllString(body, `"$1":"***"`)
	}
	return body
}