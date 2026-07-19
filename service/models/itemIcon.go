package models

import (
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

type ItemIcon struct {
	BaseModel
	IconJson        string                    `gorm:"type:varchar(1000)" json:"-"`
	Icon            datatype.ItemIconIconInfo `gorm:"-" json:"icon"`
	Title           string                    `gorm:"type:varchar(50)" json:"title"`
	Url             string                    `gorm:"type:varchar(1000)" json:"url"`
	LanUrl          string                    `gorm:"type:varchar(1000)" json:"lanUrl"`
	Description     string                    `gorm:"type:varchar(1000)" json:"description"`
	OpenMethod      int                       `gorm:"type:tinyint(1)" json:"openMethod"`
	Sort            int                       `gorm:"type:int(11)" json:"sort"`
	ItemIconGroupId int                       `json:"itemIconGroupId"`
	UserId          uint                      `json:"userId"`
	User            User                      `json:"user"`
	DepartmentId    uint                      `gorm:"default:0;index" json:"departmentId"` // 部门ID，0表示仅个人
	HealthStatus    int                       `gorm:"type:tinyint(1);default:0" json:"healthStatus"` // 0未知 1正常 2异常
	HealthCheckAt   *string                   `gorm:"type:varchar(30)" json:"healthCheckAt"`      // 最近检查时间
	HealthFailCount int                       `gorm:"type:int;default:0" json:"healthFailCount"`  // 连续失败次数
	HealthMuted     int                       `gorm:"type:tinyint(1);default:0" json:"healthMuted"` // 1维护模式(静默)
}

func (m *ItemIcon) DeleteByItemIconGroupIds(db *gorm.DB, userId uint, itemIconGroupIds []uint) (err error) {
	err = db.Delete(&ItemIcon{}, "item_icon_group_id in ? AND user_id=?", itemIconGroupIds, userId).Error
	return
}

func (m *ItemIcon) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	return db.Delete(&ItemIcon{}, "user_id=?", userId).Error
}
