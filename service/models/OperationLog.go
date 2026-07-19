package models

import "time"

// 操作日志
type OperationLog struct {
	BaseModel
	UserId      uint   `gorm:"index" json:"userId"`
	Username    string `gorm:"type:varchar(50)" json:"username"`
	Module      string `gorm:"type:varchar(50);index" json:"module"`       // 操作模块
	Action      string `gorm:"type:varchar(50)" json:"action"`             // 操作类型 create/update/delete/login
	Method      string `gorm:"type:varchar(10)" json:"method"`             // HTTP方法
	Path        string `gorm:"type:varchar(255)" json:"path"`              // 请求路径
	Ip          string `gorm:"type:varchar(50)" json:"ip"`                 // 请求IP
	UserAgent   string `gorm:"type:varchar(500)" json:"userAgent"`         // 浏览器UA
	RequestBody string `gorm:"type:text" json:"requestBody"`               // 请求参数(脱敏)
	ResponseCode int  `gorm:"type:int" json:"responseCode"`                // 响应状态码
	Duration    int64  `gorm:"type:bigint" json:"duration"`                // 耗时(ms)
	Remark      string `gorm:"type:varchar(500)" json:"remark"`            // 备注
}

func (m *OperationLog) GetList(page, pageSize int, keyword, module string, startTime, endTime *time.Time) ([]OperationLog, int64, error) {
	var list []OperationLog
	var count int64
	offset := (page - 1) * pageSize

	db := Db.Model(m)
	if keyword != "" {
		db = db.Where("username LIKE ? OR path LIKE ? OR ip LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if module != "" {
		db = db.Where("module=?", module)
	}
	if startTime != nil {
		db = db.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		db = db.Where("created_at <= ?", *endTime)
	}

	err := db.Count(&count).Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
	return list, count, err
}

func (m *OperationLog) Create() error {
	return Db.Create(m).Error
}

func (m *OperationLog) ClearAll() error {
	return Db.Where("1=1").Delete(&OperationLog{}).Error
}