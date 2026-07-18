package models

import (
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

func (m *Role) GetList(page, pageSize int) ([]Role, int64, error) {
	var list []Role
	var count int64
	offset := (page - 1) * pageSize
	err := Db.Model(m).Count(&count).Offset(offset).Limit(pageSize).Order("id asc").Find(&list).Error
	return list, count, err
}

func (m *Role) GetById(id uint) (Role, error) {
	var role Role
	err := Db.Where("id=?", id).First(&role).Error
	return role, err
}

func (m *Role) Create() (Role, error) {
	var exists Role
	if err := Db.Where("name=?", m.Name).First(&exists).Error; err == nil {
		return *m, errors.New("角色名称已存在")
	}
	err := Db.Create(m).Error
	return *m, err
}

func (m *Role) Update(id uint) error {
	var exists Role
	if err := Db.Where("name=? AND id!=?", m.Name, id).First(&exists).Error; err == nil {
		return errors.New("角色名称已存在")
	}
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
	}).Error
}

func (m *Role) UpdateFields(id uint, fields map[string]interface{}) error {
	// 如果有 name 字段，检查唯一性
	if name, ok := fields["name"]; ok {
		var exists Role
		if err := Db.Where("name=? AND id!=?", name, id).First(&exists).Error; err == nil {
			return errors.New("角色名称已存在")
		}
	}
	return Db.Model(m).Where("id=?", id).Updates(fields).Error
}

func (m *Role) Delete(id uint) error {
	var userCount int64
	Db.Model(&User{}).Where("role=?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该角色已分配用户，无法删除")
	}
	return Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&RolePermission{}, "role_id=?", id).Error; err != nil {
			return err
		}
		return tx.Delete(m, id).Error
	})
}

func (m *Role) CheckNameExists(name string, excludeId uint) bool {
	var count int64
	db := Db.Model(m).Where("name=?", name)
	if excludeId > 0 {
		db = db.Where("id!=?", excludeId)
	}
	db.Count(&count)
	return count > 0
}