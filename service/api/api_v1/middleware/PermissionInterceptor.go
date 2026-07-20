package middleware

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/lib/department"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

// PermissionInterceptor 按细粒度权限项校验访问控制。
// requiredPerm 为空字符串时表示仅需登录（个人数据接口）。
// 超级管理员(Role=1) 始终放行；其他角色按其拥有的权限项集合比对 requiredPerm。
// 部门管理员(Role=3) 也走此逻辑，按被分配的权限项生效（使部门管理员角色真正可用）。
func PermissionInterceptor(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exist := base.GetCurrentUserInfo(c)
		if !exist {
			apiReturn.ErrorByCode(c, 1001)
			c.Abort()
			return
		}

		// 超级管理员拥有全部权限，直接放行
		if currentUser.Role == department.ROLE_SUPER_ADMIN {
			c.Next()
			return
		}

		// 仅需登录即可访问的接口（个人数据域）
		if requiredPerm == "" {
			c.Next()
			return
		}

		// 其他角色：按权限项集合校验
		mRolePermission := models.RolePermission{}
		permStrs, err := mRolePermission.GetPermissionIdStringsByRoleId(currentUser.Role)
		if err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			c.Abort()
			return
		}
		for _, p := range permStrs {
			if p == requiredPerm {
				c.Next()
				return
			}
		}
		apiReturn.ErrorNoAccess(c)
		c.Abort()
	}
}

// PermissionInterceptorAny 按多个权限项的 OR 逻辑校验。
// 只要用户拥有任一 requiredPerm 即放行。适用于同一端点兼管创建/编辑等场景。
func PermissionInterceptorAny(requiredPerms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exist := base.GetCurrentUserInfo(c)
		if !exist {
			apiReturn.ErrorByCode(c, 1001)
			c.Abort()
			return
		}
		if currentUser.Role == department.ROLE_SUPER_ADMIN {
			c.Next()
			return
		}
		if len(requiredPerms) == 0 {
			c.Next()
			return
		}
		mRolePermission := models.RolePermission{}
		permStrs, err := mRolePermission.GetPermissionIdStringsByRoleId(currentUser.Role)
		if err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			c.Abort()
			return
		}
		userPermSet := make(map[string]bool, len(permStrs))
		for _, p := range permStrs {
			userPermSet[p] = true
		}
		for _, req := range requiredPerms {
			if userPermSet[req] {
				c.Next()
				return
			}
		}
		apiReturn.ErrorNoAccess(c)
		c.Abort()
	}
}
