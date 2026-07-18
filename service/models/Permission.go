package models

import "errors"

type Permission struct {
	BaseModel
	ModuleCode   string `gorm:"type:varchar(50);index" json:"moduleCode"`
	ModuleName   string `gorm:"type:varchar(100)" json:"moduleName"`
	PermissionId string `gorm:"type:varchar(50);uniqueIndex" json:"permissionId"`
	PermissionName string `gorm:"type:varchar(100)" json:"permissionName"`
	ParentId     uint   `gorm:"default:0" json:"parentId"`
	Sort         int    `gorm:"default:0" json:"sort"`
	Status       int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

func (m *Permission) GetList() ([]Permission, error) {
	var list []Permission
	err := Db.Order("sort asc").Find(&list).Error
	return list, err
}

func (m *Permission) GetByModule(moduleCode string) ([]Permission, error) {
	var list []Permission
	err := Db.Where("module_code=?", moduleCode).Order("sort asc").Find(&list).Error
	return list, err
}

func (m *Permission) GetById(id uint) (Permission, error) {
	var permission Permission
	err := Db.Where("id=?", id).First(&permission).Error
	return permission, err
}

func (m *Permission) GetByPermissionId(permissionId string) (Permission, error) {
	var permission Permission
	err := Db.Where("permission_id=?", permissionId).First(&permission).Error
	return permission, err
}

func (m *Permission) Create() (Permission, error) {
	var exists Permission
	if err := Db.Where("permission_id=?", m.PermissionId).First(&exists).Error; err == nil {
		return *m, errors.New("权限标识已存在")
	}
	err := Db.Create(m).Error
	return *m, err
}

func (m *Permission) BatchCreate(permissions []Permission) error {
	return Db.Create(&permissions).Error
}

func (m *Permission) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"module_code":   m.ModuleCode,
		"module_name":   m.ModuleName,
		"permission_id": m.PermissionId,
		"permission_name": m.PermissionName,
		"parent_id":     m.ParentId,
		"sort":          m.Sort,
		"status":        m.Status,
	}).Error
}