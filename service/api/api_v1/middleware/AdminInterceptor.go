package middleware

import (
	// "calendar-note-gin/api/v1/common/apiReturn"
	// . "calendar-note-gin/api/v1/common/base"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/lib/department"

	"github.com/gin-gonic/gin"
)

func AdminInterceptor(c *gin.Context) {
	currentUser, _ := base.GetCurrentUserInfo(c)
	if currentUser.Role != department.ROLE_SUPER_ADMIN {
		apiReturn.ErrorNoAccess(c)
		c.Abort()
		return
	}
}
