package models

import "gorm.io/gorm"

// NoticeRead 记录用户对公告/站内信的已读状态（服务端持久化，跨设备/清缓存同步）
type NoticeRead struct {
	BaseModel
	UserId   uint `gorm:"uniqueIndex:uk_notice_read;not null" json:"userId"`
	NoticeId uint `gorm:"uniqueIndex:uk_notice_read;not null" json:"noticeId"`
}

func (NoticeRead) TableName() string {
	return "notice_read"
}

// MarkRead 记录某用户已读某公告（幂等，重复调用安全）
func (m *NoticeRead) MarkRead(userId, noticeId uint) error {
	var existing NoticeRead
	err := Db.Where("user_id=? AND notice_id=?", userId, noticeId).First(&existing).Error
	if err == nil {
		return nil // 已存在，直接返回
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return Db.Create(&NoticeRead{UserId: userId, NoticeId: noticeId}).Error
}

// IsRead 判断某用户是否已读指定公告
func (m *NoticeRead) IsRead(userId, noticeId uint) (bool, error) {
	var count int64
	if err := Db.Model(&NoticeRead{}).Where("user_id=? AND notice_id=?", userId, noticeId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
