package models

import (
	"gorm.io/gorm"
)

type ItemIconGroup struct {
	BaseModel
	Icon        string `json:"icon"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(1000)" json:"description"`
	Sort         int    `gorm:"type:int(11)" json:"sort"`
	UserId       uint   `json:"userId"`
	User         User   `json:"user"`
	DepartmentId uint   `gorm:"default:0;index" json:"departmentId"` // 部门ID，0表示仅个人
}

func (m *ItemIconGroup) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	err = db.Delete(&ItemIconGroup{}, "user_id = ?", userId).Error
	return
}
