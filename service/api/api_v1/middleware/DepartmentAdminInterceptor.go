package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/lib/department"

	"github.com/gin-gonic/gin"
)

// DepartmentAdminInterceptor 部门管理员拦截器
// 超级管理员(Role=1) 和 部门管理员(Role=3) 可通过
// 其他角色拒绝访问
func DepartmentAdminInterceptor(c *gin.Context) {
	currentUser, _ := base.GetCurrentUserInfo(c)
	if currentUser.Role != department.ROLE_SUPER_ADMIN && currentUser.Role != department.ROLE_DEPT_ADMIN {
		apiReturn.ErrorNoAccess(c)
		c.Abort()
		return
	}
}