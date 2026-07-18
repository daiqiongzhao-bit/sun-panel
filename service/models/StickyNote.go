package models

import (
	"time"

	"gorm.io/gorm"
)

// StickyNote 便签
type StickyNote struct {
	BaseModel
	UserId   uint   `gorm:"index" json:"userId"`
	Content  string `gorm:"type:text" json:"content"`
	Color    string `gorm:"type:varchar(30);default:'#fff3bf'" json:"color"` // 马卡龙色
	PosX     int    `gorm:"default:100" json:"posX"`                       // X坐标
	PosY     int    `gorm:"default:100" json:"posY"`                       // Y坐标
	Width    int    `gorm:"default:200" json:"width"`                      // 宽度
	Height   int    `gorm:"default:150" json:"height"`                     // 高度
	ZIndex   int    `gorm:"default:1" json:"zIndex"`                       // 层级
	Status   int    `gorm:"type:tinyint(1);default:1" json:"status"`        // 1正常 0删除
}

func (m *StickyNote) GetByUser(userId uint) ([]StickyNote, error) {
	var list []StickyNote
	err := Db.Where("user_id=? AND status=1", userId).Order("z_index asc, updated_at desc").Find(&list).Error
	return list, err
}

func (m *StickyNote) Create() (StickyNote, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *StickyNote) Update(id, userId uint) error {
	return Db.Model(m).Where("id=? AND user_id=?", id, userId).Updates(map[string]interface{}{
		"content":  m.Content,
		"color":    m.Color,
		"pos_x":    m.PosX,
		"pos_y":    m.PosY,
		"width":    m.Width,
		"height":   m.Height,
		"z_index":  m.ZIndex,
	}).Error
}

func (m *StickyNote) UpdateWithMap(id, userId uint, updates map[string]interface{}) error {
	return Db.Model(m).Where("id=? AND user_id=?", id, userId).Updates(updates).Error
}

func (m *StickyNote) Delete(id, userId uint) error {
	// 软删除
	return Db.Model(m).Where("id=? AND user_id=?", id, userId).Update("status", 0).Error
}

// PasteBin 文本/文件中转站
type PasteBin struct {
	BaseModel
	UserId       uint      `gorm:"index" json:"userId"`
	Type         int       `gorm:"type:tinyint(1);default:1" json:"type"`             // 1文本 2文件
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	Content      string    `gorm:"type:text" json:"content"`                            // 文本内容或文件路径
	FileName     string    `gorm:"type:varchar(255)" json:"fileName"`                   // 文件名(仅type=2)
	FileSize     int64     `gorm:"default:0" json:"fileSize"`                           // 文件大小
	Code         string    `gorm:"type:varchar(20);uniqueIndex" json:"code"`            // 访问码
	Password     string    `gorm:"type:varchar(100)" json:"-"`                         // 访问密码(空=无密码, json:"-" 不返回给前端)
	BurnAfterRead int     `gorm:"type:tinyint(1);default:0" json:"burnAfterRead"`     // 1阅后即焚
	ExpireAt     time.Time `gorm:"index" json:"expireAt"`                              // 过期时间
	AccessCnt    int       `gorm:"default:0" json:"accessCnt"`                         // 访问次数
	Status       int       `gorm:"type:tinyint(1);default:1" json:"status"`             // 1有效 0过期/删除
}

func (m *PasteBin) Create() (PasteBin, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *PasteBin) GetByCode(code string) (PasteBin, error) {
	var item PasteBin
	err := Db.Where("code=? AND status=1", code).First(&item).Error
	return item, err
}

func (m *PasteBin) GetByUser(userId uint, page, pageSize int) ([]PasteBin, int64, error) {
	var list []PasteBin
	var count int64
	offset := (page - 1) * pageSize
	db := Db.Model(m).Where("user_id=?", userId)
	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("created_at desc").Find(&list).Error
	return list, count, err
}

func (m *PasteBin) IncrementAccess(id uint) error {
	return Db.Model(m).Where("id=?", id).UpdateColumn("access_cnt", gorm.Expr("access_cnt + 1")).Error
}

func (m *PasteBin) Delete(id, userId uint) error {
	return Db.Model(m).Where("id=? AND user_id=?", id, userId).Update("status", 0).Error
}

func (m *PasteBin) Update(id, userId uint, updates map[string]interface{}) error {
	return Db.Model(m).Where("id=? AND user_id=?", id, userId).Updates(updates).Error
}

// CleanExpired 清理过期的中转记录
func (m *PasteBin) CleanExpired() (int64, error) {
	result := Db.Where("status=1 AND expire_at < ?", time.Now()).Delete(m)
	return result.RowsAffected, result.Error
}