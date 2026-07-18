package models

import "errors"

type Department struct {
	BaseModel
	Name        string `gorm:"type:varchar(100)" json:"name"`
	Code        string `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	ParentId    uint   `gorm:"default:0;index" json:"parentId"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	LeaderId    uint   `gorm:"default:0" json:"leaderId"`
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"`
	Sort        int    `gorm:"default:0" json:"sort"`
}

func (m *Department) GetList() ([]Department, error) {
	var list []Department
	err := Db.Order("sort asc").Find(&list).Error
	return list, err
}

func (m *Department) GetTreeList() ([]Department, error) {
	var list []Department
	err := Db.Where("status=1").Order("sort asc").Find(&list).Error
	return list, err
}

func (m *Department) GetById(id uint) (Department, error) {
	var dept Department
	err := Db.Where("id=?", id).First(&dept).Error
	return dept, err
}

func (m *Department) Create() (Department, error) {
	var exists Department
	if err := Db.Where("code=?", m.Code).First(&exists).Error; err == nil {
		return *m, errors.New("部门编码已存在")
	}
	err := Db.Create(m).Error
	return *m, err
}

func (m *Department) Update(id uint) error {
	var exists Department
	if err := Db.Where("code=? AND id!=?", m.Code, id).First(&exists).Error; err == nil {
		return errors.New("部门编码已存在")
	}
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":        m.Name,
		"code":        m.Code,
		"parent_id":   m.ParentId,
		"description": m.Description,
		"leader_id":   m.LeaderId,
		"status":      m.Status,
		"sort":        m.Sort,
	}).Error
}

func (m *Department) Delete(id uint) error {
	// 检查子部门
	var childCount int64
	Db.Model(m).Where("parent_id=?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("该部门下存在子部门，无法删除")
	}
	// 检查是否有用户绑定
	var userCount int64
	Db.Model(&User{}).Where("department_id=?", id).Count(&userCount)
	if userCount > 0 {
		return errors.New("该部门下存在关联用户，请先解除用户绑定")
	}
	// 检查是否有图标分组绑定
	var groupCount int64
	Db.Model(&ItemIconGroup{}).Where("department_id=?", id).Count(&groupCount)
	if groupCount > 0 {
		return errors.New("该部门下存在关联的图标分组，请先移动或删除")
	}
	return Db.Delete(m, id).Error
}

func (m *Department) GetChildren(parentId uint) ([]Department, error) {
	var list []Department
	err := Db.Where("parent_id=?", parentId).Order("sort asc").Find(&list).Error
	return list, err
}