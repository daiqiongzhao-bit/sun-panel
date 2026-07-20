package models

import (
	"errors"

	"gorm.io/gorm"
)

// 用户表
type User struct {
	BaseModel
	Username     string `gorm:"index:;index:idx_username_password,priority:1;type:varchar(50)" json:"username" validate:"required"` // 账号
	Password     string `gorm:"index:idx_username_password;type:varchar(32)" json:"password" validate:"required"`                   // 密码
	Name         string `gorm:"type:varchar(20)" json:"name"`                                                                       // 名称
	HeadImage    string `gorm:"type:varchar(200)" json:"headImage"`                                                                 // 头像地址
	Status       int    `gorm:"type:tinyint(1)" json:"status"`                                                                      // 状态 1.启用 2.停用 3.未激活
	Role         int    `gorm:"type:int(11)" json:"role"`                                                                           // 角色 1.管理员 2.普通用户
	Mail         string `gorm:"type:varchar(50)" json:"mail"`                                                                       // 邮箱
	ReferralCode string `gorm:"type:varchar(10)" json:"referralCode"`                                                               // 推荐码
	Token        string `gorm:"type:varchar(32)" json:"token"`
	DepartmentId uint   `gorm:"default:0;index" json:"departmentId"`                                                                // 部门ID，0表示无部门

	// 登录安全
	LoginFailCount int   `gorm:"type:int;default:0" json:"loginFailCount"`   // 连续登录失败次数
	LoginLockUntil int64 `gorm:"type:bigint;default:0" json:"loginLockUntil"` // 锁定截止时间(unix秒)，0=未锁定
	TwoFAEnabled   int   `gorm:"type:tinyint(1);default:0" json:"twoFAEnabled"` // 是否启用两步验证 1=启用
	TwoFASecret    string `gorm:"type:varchar(64)" json:"-"`                 // TOTP 密钥(不对外暴露)

	Permissions  []string `gorm:"-" json:"permissions"` // 当前用户拥有的权限标识列表（运行时填充，不入库）
	UserId uint `gorm:"-"  json:"userId"`
}

// 获取用户信息
func (m *User) GetUserInfoByUid(uid uint) (User, error) {
	mUser := User{}
	err := Db.Where("id=?", uid).First(&mUser).Error
	return mUser, err
}

// 根据用户名和密码查询用户
func (m *User) GetUserInfoByUsernameAndPassword(username, password string) (User, error) {
	userInfo := User{}
	err := Db.Where("username=?", username).Where("password=?", password).First(&userInfo).Error
	return userInfo, err
}

// 根据用户名查询用户
func (m *User) GetUserInfoByUsername(username string) (User, error) {
	mUser := User{}
	err := Db.Where("username=?", username).First(&mUser).Error
	return mUser, err
}

// 根据邮箱查询用户
func (m *User) GetUserInfoByMail() *User {
	mUser := User{}
	if Db.Where("mail=?", m.Mail).First(&mUser).Error != nil {
		return nil
	}
	return &mUser
}

// 根据token查询用户
func (m *User) GetUserInfoByToken(userToken string) (User, error) {
	mUser := User{}
	err := Db.Where("token=?", userToken).First(&mUser).Error
	return mUser, err
}

// 更新用户基于id
// 支持：name,autograph,header_image,status,role,mail,token,password,username,gender
func (m *User) UpdateUserInfoByUserId(user_id uint, updateInfo map[string]interface{}) error {
	mUser := User{}

	data := map[string]interface{}{}
	if v, ok := updateInfo["name"]; ok {
		data["name"] = v
	}
	if v, ok := updateInfo["head_image"]; ok {
		data["head_image"] = v
	}
	if v, ok := updateInfo["status"]; ok {
		data["status"] = v
	}
	if v, ok := updateInfo["role"]; ok {
		data["role"] = v
	}
	if v, ok := updateInfo["department_id"]; ok {
		data["department_id"] = v
	}
	if v, ok := updateInfo["gender"]; ok {
		data["gender"] = v
	}

	if v, ok := updateInfo["mail"]; ok {
		hasUser := User{}
		count := Db.Where("mail=?", updateInfo["mail"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != user_id {
			return errors.New("the mail already exists")
		}
		data["mail"] = v
	}
	if v, ok := updateInfo["username"]; ok {
		hasUser := User{}
		count := Db.Where("username=?", updateInfo["username"]).First(&hasUser).RowsAffected
		if count != 0 && hasUser.ID != user_id {
			return errors.New("the username already exists")
		}
		data["username"] = v
	}
	if v, ok := updateInfo["token"]; ok {
		data["token"] = v
	}
	if v, ok := updateInfo["password"]; ok {
		data["password"] = v
	}

	err := Db.Model(&mUser).Where("id=?", user_id).Updates(data).Error

	return err
}

// 添加一个
func (m *User) CreateOne() (User, error) {
	err := Db.Create(m).Error
	return *m, err
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckMailAndUsername(mail, username string) error {
	hasUser := User{}
	count := Db.Where("mail=?", mail).First(&hasUser).RowsAffected
	if count != 0 {
		return errors.New("该邮箱已被注册")
	}

	count = Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return errors.New("该用户名已被注册")
	}
	return nil
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckMailExist(mail string) (User, error) {
	hasUser := User{}
	count := Db.Where("mail=?", mail).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该邮箱已被注册")
	}
	return hasUser, nil
}

// 验证是否有重复的用户名或者邮箱
func (m *User) CheckUsernameExist(username string) (User, error) {
	hasUser := User{}
	count := Db.Where("username=?", username).First(&hasUser).RowsAffected
	if count != 0 {
		return hasUser, errors.New("该用户名已被注册")
	}
	return hasUser, nil
}

// // 根据用户名和密码查询用户
// func (m *User) CreateUser(uid uint) *User {
// 	mUser := User{}
// 	if Db.Where("id=?", uid).First(&mUser).Error != nil {
// 		return nil
// 	} else {
// 		return &mUser
// 	}
// }

// UpdateTwoFA 更新两步验证状态与密钥
func (m *User) UpdateTwoFA(userId uint, enabled int, secret string) error {
	data := map[string]interface{}{"two_fa_enabled": enabled}
	if secret != "" {
		data["two_fa_secret"] = secret
	}
	return Db.Model(&User{}).Where("id=?", userId).Updates(data).Error
}

// IncrementLoginFail 登录失败计数 +1，并在达到阈值时锁定至 lockUntil
func (m *User) IncrementLoginFail(userId uint, lockUntil int64) error {
	return Db.Model(&User{}).Where("id=?", userId).Updates(map[string]interface{}{
		"login_fail_count": gorm.Expr("login_fail_count + 1"),
		"login_lock_until": lockUntil,
	}).Error
}

// ResetLoginFail 登录成功后重置失败计数与锁定
func (m *User) ResetLoginFail(userId uint) error {
	return Db.Model(&User{}).Where("id=?", userId).Updates(map[string]interface{}{
		"login_fail_count": 0,
		"login_lock_until": 0,
	}).Error
}
