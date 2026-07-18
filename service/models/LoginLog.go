package models

// 登录日志
type LoginLog struct {
	BaseModel
	UserId    uint   `gorm:"index" json:"userId"`
	Username  string `gorm:"type:varchar(50)" json:"username"`
	Ip        string `gorm:"type:varchar(50)" json:"ip"`               // 登录IP
	Location  string `gorm:"type:varchar(100)" json:"location"`        // IP归属地(预留)
	UserAgent string `gorm:"type:varchar(500)" json:"userAgent"`       // 浏览器UA
	Status    int    `gorm:"type:tinyint(1)" json:"status"`             // 登录状态 1成功 2失败
	Remark    string `gorm:"type:varchar(255)" json:"remark"`           // 备注(失败原因等)
}

func (m *LoginLog) GetList(page, pageSize int, keyword string) ([]LoginLog, int64, error) {
	var list []LoginLog
	var count int64
	offset := (page - 1) * pageSize

	db := Db.Model(m)
	if keyword != "" {
		db = db.Where("username LIKE ? OR ip LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
	return list, count, err
}

func (m *LoginLog) Create() error {
	return Db.Create(m).Error
}

func (m *LoginLog) ClearAll() error {
	return Db.Where("1=1").Delete(&LoginLog{}).Error
}