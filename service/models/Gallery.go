package models

import "errors"

type Gallery struct {
	BaseModel
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Src         string `gorm:"type:varchar(500)" json:"src"`
	CategoryId  uint   `gorm:"default:0;index" json:"categoryId"`
	GalleryType int    `gorm:"type:tinyint(1);default:1" json:"galleryType"`
	FileSize    int64  `json:"fileSize"`
	UserId      uint   `gorm:"index" json:"userId"`
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"`
}

type GalleryCategory struct {
	BaseModel
	Name        string `gorm:"type:varchar(100)" json:"name"`
	Code        string `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	GalleryType int    `gorm:"type:tinyint(1)" json:"galleryType"`
	Sort        int    `gorm:"default:0" json:"sort"`
}

func (m *Gallery) GetList(categoryId uint, galleryType int, keyword string, page, pageSize int) ([]Gallery, int64, error) {
	var list []Gallery
	var count int64
	db := Db.Model(m).Where("status=1")
	if categoryId > 0 {
		db = db.Where("category_id=?", categoryId)
	}
	if galleryType > 0 {
		db = db.Where("gallery_type=?", galleryType)
	}
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	offset := (page - 1) * pageSize
	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error
	return list, count, err
}

func (m *Gallery) GetById(id uint) (Gallery, error) {
	var gallery Gallery
	err := Db.Where("id=?", id).First(&gallery).Error
	return gallery, err
}

func (m *Gallery) Create() (Gallery, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *Gallery) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":         m.Name,
		"src":          m.Src,
		"category_id":  m.CategoryId,
		"gallery_type": m.GalleryType,
		"status":       m.Status,
	}).Error
}

func (m *Gallery) Delete(id uint) error {
	return Db.Delete(m, id).Error
}

func (m *Gallery) BatchDelete(ids []uint) error {
	return Db.Delete(m, "id IN ?", ids).Error
}

func (m *GalleryCategory) GetList(galleryType int) ([]GalleryCategory, error) {
	var list []GalleryCategory
	db := Db.Model(m)
	if galleryType > 0 {
		db = db.Where("gallery_type=?", galleryType)
	}
	err := db.Order("sort asc").Find(&list).Error
	return list, err
}

func (m *GalleryCategory) GetById(id uint) (GalleryCategory, error) {
	var category GalleryCategory
	err := Db.Where("id=?", id).First(&category).Error
	return category, err
}

func (m *GalleryCategory) Create() (GalleryCategory, error) {
	var exists GalleryCategory
	if err := Db.Where("code=?", m.Code).First(&exists).Error; err == nil {
		return *m, errors.New("分类编码已存在")
	}
	err := Db.Create(m).Error
	return *m, err
}

func (m *GalleryCategory) Update(id uint) error {
	var exists GalleryCategory
	if err := Db.Where("code=? AND id!=?", m.Code, id).First(&exists).Error; err == nil {
		return errors.New("分类编码已存在")
	}
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"name":         m.Name,
		"code":         m.Code,
		"gallery_type": m.GalleryType,
		"sort":         m.Sort,
	}).Error
}

func (m *GalleryCategory) Delete(id uint) error {
	var galleryCount int64
	Db.Model(&Gallery{}).Where("category_id=?", id).Count(&galleryCount)
	if galleryCount > 0 {
		return errors.New("该分类下存在图片，无法删除")
	}
	return Db.Delete(m, id).Error
}