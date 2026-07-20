package models

import (
	"gorm.io/gorm"
)

type RolePermission struct {
	BaseModel
	RoleId       uint `gorm:"index" json:"roleId"`
	PermissionId uint `gorm:"index" json:"permissionId"`
}

func (m *RolePermission) GetByRoleId(roleId uint) ([]RolePermission, error) {
	var list []RolePermission
	err := Db.Where("role_id=?", roleId).Find(&list).Error
	return list, err
}

func (m *RolePermission) GetPermissionIdsByRoleId(roleId uint) ([]uint, error) {
	var ids []uint
	err := Db.Model(m).Where("role_id=?", roleId).Pluck("permission_id", &ids).Error
	return ids, err
}

// GetPermissionIdStringsByRoleId 返回该角色拥有的权限标识字符串列表（如 "user:delete"）
func (m *RolePermission) GetPermissionIdStringsByRoleId(roleId int) ([]string, error) {
	var permIds []uint
	if err := Db.Model(m).Where("role_id=?", uint(roleId)).Pluck("permission_id", &permIds).Error; err != nil {
		return nil, err
	}
	if len(permIds) == 0 {
		return []string{}, nil
	}
	var permStrs []string
	if err := Db.Model(&Permission{}).Where("id IN ?", permIds).Pluck("permission_id", &permStrs).Error; err != nil {
		return nil, err
	}
	return permStrs, nil
}

func (m *RolePermission) BatchSave(roleId uint, permissionIds []uint) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&RolePermission{}, "role_id=?", roleId).Error; err != nil {
			return err
		}
		if len(permissionIds) == 0 {
			return nil
		}
		var records []RolePermission
		for _, pid := range permissionIds {
			records = append(records, RolePermission{
				RoleId:       roleId,
				PermissionId: pid,
			})
		}
		return tx.Create(&records).Error
	})
}

func (m *RolePermission) DeleteByRoleId(roleId uint) error {
	return Db.Delete(&RolePermission{}, "role_id=?", roleId).Error
}

func (m *RolePermission) DeleteByPermissionId(permissionId uint) error {
	return Db.Delete(&RolePermission{}, "permission_id=?", permissionId).Error
}