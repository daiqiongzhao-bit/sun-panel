package department

import (
	"sun-panel/models"

	"gorm.io/gorm"
)

// 角色常量
const (
	ROLE_SUPER_ADMIN = 1 // 超级管理员（原 Role=1）
	ROLE_DEPT_ADMIN  = 3 // 部门管理员
)

// GetUserVisibleDepartmentIds 获取用户可见的部门ID列表
// 返回值：部门ID列表 + 是否为超级管理员
// 超级管理员：返回nil（表示不限制，可看所有数据）
// 部门管理员：返回本部门 + 所有子部门的ID
// 普通用户：返回仅本部门的ID
func GetUserVisibleDepartmentIds(user models.User) ([]uint, bool) {
	// 超级管理员不受限制
	if user.Role == ROLE_SUPER_ADMIN {
		return nil, true
	}

	if user.DepartmentId == 0 {
		// 用户未分配部门：严格只看自己的数据，不做任何部门共享
		return []uint{}, false
	}

	// 部门管理员：获取本部门 + 所有子部门
	if user.Role == ROLE_DEPT_ADMIN {
		ids := []uint{user.DepartmentId}
		ids = append(ids, getChildDepartmentIds(user.DepartmentId)...)
		return ids, false
	}

	// 普通用户：只能看本部门的数据
	return []uint{user.DepartmentId}, false
}

// getChildDepartmentIds 递归获取所有子部门ID
func getChildDepartmentIds(parentId uint) []uint {
	var ids []uint
	children, err := (&models.Department{}).GetChildren(parentId)
	if err != nil {
		return ids
	}
	for _, child := range children {
		ids = append(ids, child.ID)
		ids = append(ids, getChildDepartmentIds(child.ID)...)
	}
	return ids
}

// BuildDepartmentScope 构建 GORM 的部门过滤条件
// 超级管理员：不加任何过滤
// 其他角色：WHERE (user_id = 当前用户 OR department_id IN (可见部门列表))
// 这样用户既能看自己的数据，也能看同部门其他人分享到部门的数据
func BuildDepartmentScope(db *gorm.DB, userId uint, user models.User) *gorm.DB {
	deptIds, isSuperAdmin := GetUserVisibleDepartmentIds(user)
	if isSuperAdmin {
		return db
	}

	// 用户可以看到：自己的数据 + 所在部门的数据
	if len(deptIds) == 0 {
		// 无部门或部门范围为空：仅过滤为自己的数据
		return db.Where("user_id = ?", userId)
	}
	return db.Where("user_id = ? OR department_id IN ?", userId, deptIds)
}