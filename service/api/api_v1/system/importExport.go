package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ImportExportApi struct{}

// ExportUsers 导出所有用户
func (a *ImportExportApi) ExportUsers(c *gin.Context) {
	var users []models.User
	if err := global.Db.Find(&users).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 脱敏：清除密码
	for i := range users {
		users[i].Password = ""
		users[i].Token = ""
	}

	apiReturn.SuccessData(c, users)
}

// ImportUsers 导入用户
func (a *ImportExportApi) ImportUsers(c *gin.Context) {
	type Request struct {
		Users []models.User `json:"users"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	imported := 0
	skipped := 0
	for _, user := range req.Users {
		// 检查用户名是否已存在
		var existing models.User
		if err := global.Db.Where("username=?", user.Username).First(&existing).Error; err == nil {
			skipped++
			continue
		}
		user.Password = "" // 导入的用户需要重置密码
		if err := global.Db.Create(&user).Error; err == nil {
			imported++
		}
	}

	apiReturn.SuccessData(c, gin.H{"imported": imported, "skipped": skipped})
}

// ExportRoles 导出角色及权限
func (a *ImportExportApi) ExportRoles(c *gin.Context) {
	var roles []models.Role
	if err := global.Db.Find(&roles).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	type RoleExport struct {
		Role        models.Role          `json:"role"`
		Permissions []models.Permission  `json:"permissions"`
	}

	result := []RoleExport{}
	for _, role := range roles {
		var rps []models.RolePermission
		global.Db.Where("role_id=?", role.ID).Find(&rps)

		var perms []models.Permission
		for _, rp := range rps {
			var perm models.Permission
			if err := global.Db.Where("id=?", rp.PermissionId).First(&perm).Error; err == nil {
				perms = append(perms, perm)
			}
		}

		result = append(result, RoleExport{Role: role, Permissions: perms})
	}

	apiReturn.SuccessData(c, result)
}

// ImportRoles 导入角色及权限
func (a *ImportExportApi) ImportRoles(c *gin.Context) {
	type RoleImport struct {
		Role        models.Role         `json:"role"`
		Permissions []models.Permission `json:"permissions"`
	}
	type Request struct {
		Roles []RoleImport `json:"roles"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	imported := 0
	for _, ri := range req.Roles {
		// 检查角色名是否已存在
		var existing models.Role
		if err := global.Db.Where("name=?", ri.Role.Name).First(&existing).Error; err == nil {
			continue
		}

		// 创建角色
		role := models.Role{
			Name:        ri.Role.Name,
			Description: ri.Role.Description,
			Status:      ri.Role.Status,
		}
		if err := global.Db.Create(&role).Error; err != nil {
			continue
		}

		// 分配权限
		for _, perm := range ri.Permissions {
			// 查找或创建权限
			var existingPerm models.Permission
			if err := global.Db.Where("permission_name=?", perm.PermissionName).First(&existingPerm).Error; err != nil {
				if err2 := global.Db.Create(&perm).Error; err2 == nil {
					existingPerm = perm
				} else {
					continue
				}
			}
			rp := models.RolePermission{RoleId: role.ID, PermissionId: existingPerm.ID}
			global.Db.Create(&rp)
		}
		imported++
	}

	apiReturn.SuccessData(c, gin.H{"imported": imported})
}

