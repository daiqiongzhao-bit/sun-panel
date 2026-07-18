package models

import (
	"errors"

	"gorm.io/gorm"
)

type Icon struct {
	BaseModel
	Name        string `gorm:"type:varchar(100)" json:"name"`
	IconKey     string `gorm:"type:varchar(100);index" json:"iconKey"`
	CategoryId  uint   `gorm:"default:0;index" json:"categoryId"`
	Src         string `gorm:"type:varchar(500)" json:"src"`
	IconType    int    `gorm:"type:tinyint(1);default:1" json:"iconType"`
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

type IconCategory struct {
	BaseModel
	Name        string `gorm:"type:varchar(100)" json:"name"`
	Code        string `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	ParentId    uint   `gorm:"default:0" json:"parentId"`
	Sort        int    `gorm:"default:0" json:"sort"`
}

type IconFavorite struct {
	BaseModel
	UserId uint `gorm:"index" json:"userId"`
	IconId uint `gorm:"index" json:"iconId"`
}

func (m *Icon) GetList(categoryId uint, keyword string, page, pageSize int) ([]Icon, int64, error) {
	var list []Icon
	var count int64
	db := Db.Model(m).Where("status=1")
	if categoryId > 0 {
		db = db.Where("category_id=?", categoryId)
	}
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	offset := (page - 1) * pageSize
	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("sort asc").Find(&list).Error
	return list, count, err
}

func (m *Icon) GetById(id uint) (Icon, error) {
	var icon Icon
	err := Db.Where("id=?", id).First(&icon).Error
	return icon, err
}

func (m *Icon) GetByIconKey(iconKey string) (Icon, error) {
	var icon Icon
	err := Db.Where("icon_key=?", iconKey).First(&icon).Error
	return icon, err
}

func (m *Icon) BatchGet(ids []uint) ([]Icon, error) {
	var list []Icon
	err := Db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (m *Icon) GetByCategory(categoryId uint) ([]Icon, error) {
	var list []Icon
	err := Db.Where("category_id=? AND status=1", categoryId).Order("sort asc").Find(&list).Error
	return list, err
}

func (m *Icon) Create() (Icon, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *Icon) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":       m.Name,
		"icon_key":   m.IconKey,
		"category_id": m.CategoryId,
		"src":        m.Src,
		"icon_type":  m.IconType,
		"status":     m.Status,
	}).Error
}

func (m *Icon) Delete(id uint) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&IconFavorite{}, "icon_id=?", id).Error; err != nil {
			return err
		}
		return tx.Delete(m, id).Error
	})
}

func (m *IconCategory) GetList() ([]IconCategory, error) {
	var list []IconCategory
	err := Db.Order("sort asc").Find(&list).Error
	return list, err
}

func (m *IconCategory) GetById(id uint) (IconCategory, error) {
	var category IconCategory
	err := Db.Where("id=?", id).First(&category).Error
	return category, err
}

func (m *IconCategory) Create() (IconCategory, error) {
	var exists IconCategory
	if err := Db.Where("code=?", m.Code).First(&exists).Error; err == nil {
		return *m, errors.New("分类编码已存在")
	}
	err := Db.Create(m).Error
	return *m, err
}

func (m *IconCategory) Update(id uint) error {
	var exists IconCategory
	if err := Db.Where("code=? AND id!=?", m.Code, id).First(&exists).Error; err == nil {
		return errors.New("分类编码已存在")
	}
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":      m.Name,
		"code":      m.Code,
		"parent_id": m.ParentId,
		"sort":      m.Sort,
	}).Error
}

func (m *IconCategory) Delete(id uint) error {
	var iconCount int64
	Db.Model(&Icon{}).Where("category_id=?", id).Count(&iconCount)
	if iconCount > 0 {
		return errors.New("该分类下存在图标，无法删除")
	}
	return Db.Delete(m, id).Error
}

func (m *IconFavorite) IsFavorite(userId, iconId uint) bool {
	var count int64
	Db.Model(m).Where("user_id=? AND icon_id=?", userId, iconId).Count(&count)
	return count > 0
}

func (m *IconFavorite) Toggle(userId, iconId uint) error {
	if m.IsFavorite(userId, iconId) {
		return Db.Delete(m, "user_id=? AND icon_id=?", userId, iconId).Error
	}
	return Db.Create(&IconFavorite{UserId: userId, IconId: iconId}).Error
}

func (m *IconFavorite) GetFavorites(userId uint) ([]uint, error) {
	var ids []uint
	err := Db.Model(m).Where("user_id=?", userId).Pluck("icon_id", &ids).Error
	return ids, err
}

func (m *Icon) GetByName(name string) (Icon, error) {
	var icon Icon
	err := Db.Where("name=?", name).First(&icon).Error
	return icon, err
}

func (m *Icon) GetByIds(ids []uint) ([]Icon, error) {
	var list []Icon
	err := Db.Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (m *IconFavorite) GetByUserId(userId uint) ([]Icon, error) {
	var iconIds []uint
	err := Db.Model(m).Where("user_id=?", userId).Pluck("icon_id", &iconIds).Error
	if err != nil {
		return nil, err
	}
	var icon Icon
	return icon.GetByIds(iconIds)
}