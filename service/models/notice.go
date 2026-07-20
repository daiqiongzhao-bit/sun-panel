package models

import "fmt"

const (
	NOTICE_DISPLAY_TYPE_LOGIN = iota + 1 // 通知展示类型 登录页
	NOTICE_DISPLAY_TYPE_HOME             // 通知展示类型 首页
)

const (
	NOTICE_TYPE_ANNOUNCEMENT = 1 // 公告（全员可见）
	NOTICE_TYPE_MESSAGE      = 2 // 站内信（定向用户）
)

const (
	NOTICE_STATUS_ENABLED  = 1 // 启用
	NOTICE_STATUS_DISABLED = 2 // 停用
)

type Notice struct {
	BaseModel
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	DisplayType int    `gorm:"type:tinyint(1)" json:"displayType"` // 展示类型 参考常量
	NoticeType  int    `gorm:"type:tinyint(1);default:1" json:"noticeType"` // 1公告 2站内信
	OneRead     int    `gorm:"type:tinyint(1)" json:"oneRead"`     // 1.前端记录读取状态 0.每次都展示
	Url         string `gorm:"type:varchar(255)" json:"url"`       // 跳转地址
	IsLogin     uint   `gorm:"type:tinyint(1)" json:"isLogin"`     // 登录可见
	Status      int    `gorm:"type:tinyint(1);default:1" json:"status"` // 状态 1启用 2停用
	UserId      uint   `gorm:"index" json:"userId"`                // 发布人ID
	TargetUserIds string `gorm:"type:varchar(2000)" json:"targetUserIds"` // 站内信目标用户ID,逗号分隔,空=全员
	User        User   `json:"user"`
}

func (m *Notice) GetList(page, pageSize int, keyword string, noticeType int) ([]Notice, int64, error) {
	var list []Notice
	var count int64
	offset := (page - 1) * pageSize

	db := Db.Model(m).Where("status=?", NOTICE_STATUS_ENABLED)
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if noticeType > 0 {
		db = db.Where("notice_type=?", noticeType)
	}

	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
	return list, count, err
}

func (m *Notice) GetById(id uint) (Notice, error) {
	var notice Notice
	err := Db.Where("id=?", id).First(&notice).Error
	return notice, err
}

func (m *Notice) GetVisibleNotices(userId uint, displayTypes ...int) ([]Notice, error) {
	var list []Notice
	if len(displayTypes) == 0 {
		displayTypes = []int{NOTICE_DISPLAY_TYPE_HOME}
	}
	// 公告全员可见 + 站内信匹配目标用户
	err := Db.Where("(notice_type=? OR (notice_type=? AND (target_user_ids='' OR target_user_ids LIKE ?))) AND display_type IN ? AND status=?",
		NOTICE_TYPE_ANNOUNCEMENT, NOTICE_TYPE_MESSAGE,
		"%,"+fmt.Sprintf("%d", userId)+",%",
		displayTypes, NOTICE_STATUS_ENABLED,
	).Order("id desc").Find(&list).Error
	return list, err
}

func (m *Notice) Create() (Notice, error) {
	err := Db.Create(m).Error
	return *m, err
}

func (m *Notice) Update(id uint) error {
	return Db.Model(m).Where("id=?", id).Updates(map[string]interface{}{
		"title":         m.Title,
		"content":       m.Content,
		"display_type":  m.DisplayType,
		"notice_type":   m.NoticeType,
		"one_read":      m.OneRead,
		"url":           m.Url,
		"is_login":      m.IsLogin,
		"status":        m.Status,
		"target_user_ids": m.TargetUserIds,
	}).Error
}

func (m *Notice) Delete(id uint) error {
	return Db.Delete(m, id).Error
}