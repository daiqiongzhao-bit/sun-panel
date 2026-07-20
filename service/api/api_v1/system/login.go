package system

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/totp"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginApi struct {
}

// 登录输入验证
type LoginLoginVerify struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=50"`
	VCode    string `json:"vcode" validate:"max=6"`
	Email    string `json:"email"`
}

// 2FA 第二步校验输入
type LoginTwoFAVerify struct {
	TwoFaToken string `json:"twoFaToken" validate:"required"`
	Code       string `json:"code" validate:"required,len=6"`
}

// 2FA 自助管理-确认/禁用/验证输入
type TwoFAVerifyCodeReq struct {
	Code string `json:"code" validate:"required,len=6"`
}

// 登录安全相关配置
const (
	securityMaxFail     = 5           // 最大连续登录失败次数
	securityLockMinutes = 15          // 达到上限后的锁定分钟数
	twoFaIssuer         = "Sun-Panel" // TOTP otpauth issuer 名称
)

// twoFaPending: 登录阶段生成的临时 token -> userId（2FA 第二步校验用，登录成功后清除）
var twoFaPending sync.Map

// twoFaSetup: 启用 2FA 时临时保存 secret -> userId（confirm 通过后才落库启用）
var twoFaSetup sync.Map

// @Summary 登录账号
// @Accept application/json
// @Produce application/json
// @Param LoginLoginVerify body LoginLoginVerify true "登陆验证信息"
// @Tags user
// @Router /login [post]
// checkLoginAllowIps 校验客户端 IP 是否在允许登录的白名单内。
// 未配置（空）则放行；配置读取失败时也放行，避免把自己锁死。
func checkLoginAllowIps(c *gin.Context) bool {
	cfg := systemSetting.ApplicationSetting{}
	if err := global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &cfg); err != nil {
		return true
	}
	allow := strings.TrimSpace(cfg.Login.LoginAllowIps)
	if allow == "" {
		return true
	}
	clientIP := net.ParseIP(c.ClientIP())
	if clientIP == nil {
		// 无法解析客户端 IP 时放行，避免误杀
		return true
	}
	for _, item := range strings.Split(allow, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if strings.Contains(item, "/") {
			if _, ipnet, err := net.ParseCIDR(item); err == nil && ipnet.Contains(clientIP) {
				return true
			}
		} else if allowIP := net.ParseIP(item); allowIP != nil && allowIP.Equal(clientIP) {
			return true
		}
	}
	return false
}

func (l LoginApi) Login(c *gin.Context) {
	// IP 白名单校验（仅限制登录入口，不影响已登录访问）
	if !checkLoginAllowIps(c) {
		recordLoginLog(0, "", c, 2, "IP 不在允许登录白名单")
		apiReturn.Error(c, "当前 IP 不在允许登录的列表中")
		return
	}

	param := LoginLoginVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mUser := models.User{}
	var (
		err  error
		info models.User
	)
	param.Username = strings.TrimSpace(param.Username)

	// 先按用户名查用户（便于在密码错误时记录 userId 并做失败计数/锁定）
	if info, err = mUser.GetUserInfoByUsername(param.Username); err != nil {
		if err == gorm.ErrRecordNotFound {
			recordLoginLog(0, param.Username, c, 2, "账号不存在")
			apiReturn.ErrorByCode(c, 1003)
			return
		}
		recordLoginLog(0, param.Username, c, 2, err.Error())
		apiReturn.Error(c, err.Error())
		return
	}

	// 账号仍在锁定中：拒绝并提示剩余秒数
	if info.LoginLockUntil > time.Now().Unix() {
		remaining := info.LoginLockUntil - time.Now().Unix()
		recordLoginLog(info.ID, info.Username, c, 2, "账号已锁定")
		apiReturn.Error(c, "账号已锁定，请 "+strconv.FormatInt(remaining, 10)+" 秒后重试")
		return
	}
	// 锁定已过期：解锁并清空失败计数，重新开始计数
	if info.LoginLockUntil > 0 {
		mUser.ResetLoginFail(info.ID)
		info.LoginFailCount = 0
	}

	// 密码校验
	if info.Password != cmn.PasswordEncryption(param.Password) {
		failCount := info.LoginFailCount + 1
		lockUntil := int64(0)
		if failCount >= securityMaxFail {
			lockUntil = time.Now().Add(securityLockMinutes * time.Minute).Unix()
		}
		mUser.IncrementLoginFail(info.ID, lockUntil)
		remark := "密码错误"
		if lockUntil > 0 {
			remark = "密码错误，账号已锁定"
		}
		recordLoginLog(info.ID, info.Username, c, 2, remark)
		apiReturn.ErrorByCode(c, 1003)
		return
	}

	// 停用或未激活
	if info.Status != 1 {
		recordLoginLog(info.ID, info.Username, c, 2, "账号已停用或未激活")
		apiReturn.ErrorByCode(c, 1004)
		return
	}

	// 已启用两步验证：不直接下发最终 token，返回 needTwoFA 让前端进入第二步
	if info.TwoFAEnabled == 1 {
		twoFaToken := cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
		twoFaPending.Store(twoFaToken, info.ID)
		recordLoginLog(info.ID, info.Username, c, 1, "两步验证待校验")
		apiReturn.SuccessData(c, gin.H{
			"needTwoFA":  true,
			"twoFaToken": twoFaToken,
		})
		return
	}

	// 正常登录：下发 token 并重置失败计数
	mUser.ResetLoginFail(info.ID)
	issueToken(c, info)
}

// @Summary 两步验证(2FA)第二步
// @Accept application/json
// @Produce application/json
// @Param LoginTwoFAVerify body LoginTwoFAVerify true "两步验证信息"
// @Tags user
// @Router /login/2fa [post]
func (l LoginApi) Login2FA(c *gin.Context) {
	if !checkLoginAllowIps(c) {
		apiReturn.Error(c, "当前 IP 不在允许登录的列表中")
		return
	}

	param := LoginTwoFAVerify{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	val, ok := twoFaPending.LoadAndDelete(param.TwoFaToken)
	if !ok {
		apiReturn.Error(c, "两步验证会话已失效，请重新登录")
		return
	}
	userId := val.(uint)

	mUser := models.User{}
	info, err := mUser.GetUserInfoByUid(userId)
	if err != nil {
		apiReturn.Error(c, "用户不存在")
		return
	}
	if info.TwoFAEnabled != 1 || info.TwoFASecret == "" {
		apiReturn.Error(c, "该账号未启用两步验证")
		return
	}
	if !totp.Validate(info.TwoFASecret, param.Code) {
		recordLoginLog(info.ID, info.Username, c, 2, "两步验证码错误")
		apiReturn.Error(c, "验证码错误")
		return
	}
	mUser.ResetLoginFail(info.ID)
	issueToken(c, info)
}

// issueToken 生成并下发最终登录 token（从原 Login 抽取复用）
func issueToken(c *gin.Context, info models.User) {
	bToken := info.Token
	if info.Token == "" {
		buildTokenOver := false
		mUser := models.User{}
		for !buildTokenOver {
			bToken = cmn.BuildRandCode(32, cmn.RAND_CODE_MODE2)
			if _, err := mUser.GetUserInfoByToken(bToken); err != nil {
				// 保存token
				mUser.UpdateUserInfoByUserId(info.ID, map[string]interface{}{
					"token": bToken,
				})
				buildTokenOver = true
			}
		}
		info.Token = bToken
	}
	info.Password = ""
	info.ReferralCode = ""

	// global.UserToken.SetDefault(bToken, info)
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(info.ID))))
	global.CUserToken.SetDefault(cToken, bToken)
	global.Logger.Debug("token:", cToken, "|", bToken)
	global.Logger.Debug(global.CUserToken.Get(cToken))

	// 设置当前用户信息
	c.Set("userInfo", info)
	info.Token = cToken // 重要 采用cToken,隐藏真实token

	// 记录登录成功日志
	recordLoginLog(info.ID, info.Username, c, 1, "登录成功")

	apiReturn.SuccessData(c, info)
}

// @Summary 获取当前账号两步验证状态
// @Tags user
// @Router /twofa/status [post]
func (l LoginApi) TwoFAStatus(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	apiReturn.SuccessData(c, gin.H{"enabled": userInfo.TwoFAEnabled == 1})
}

// @Summary 生成两步验证密钥(暂不启用)，返回 otpauth URI 供前端展示二维码
// @Tags user
// @Router /twofa/enable [post]
func (l LoginApi) TwoFAEnable(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	secret, err := totp.GenerateSecret()
	if err != nil {
		apiReturn.Error(c, "生成密钥失败")
		return
	}
	twoFaSetup.Store(userInfo.ID, secret)
	uri := totp.OTPAuthURI(userInfo.Username, twoFaIssuer, secret)
	apiReturn.SuccessData(c, gin.H{
		"secret":  secret,
		"otpauth": uri,
	})
}

// @Summary 校验动态码并启用两步验证
// @Tags user
// @Router /twofa/confirm [post]
func (l LoginApi) TwoFAConfirm(c *gin.Context) {
	param := TwoFAVerifyCodeReq{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	val, ok := twoFaSetup.LoadAndDelete(userInfo.ID)
	if !ok {
		apiReturn.Error(c, "请先获取两步验证密钥")
		return
	}
	secret := val.(string)
	if !totp.Validate(secret, param.Code) {
		apiReturn.Error(c, "验证码错误，启用失败")
		return
	}
	mUser := models.User{}
	if err := mUser.UpdateTwoFA(userInfo.ID, 1, secret); err != nil {
		apiReturn.Error(c, "启用失败")
		return
	}
	apiReturn.Success(c)
}

// @Summary 校验动态码并关闭两步验证
// @Tags user
// @Router /twofa/disable [post]
func (l LoginApi) TwoFADisable(c *gin.Context) {
	param := TwoFAVerifyCodeReq{}
	if err := c.ShouldBindJSON(&param); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	mUser := models.User{}
	info, err := mUser.GetUserInfoByUid(userInfo.ID)
	if err != nil || info.TwoFASecret == "" {
		apiReturn.Error(c, "未启用两步验证")
		return
	}
	if !totp.Validate(info.TwoFASecret, param.Code) {
		apiReturn.Error(c, "验证码错误，关闭失败")
		return
	}
	if err := mUser.UpdateTwoFA(userInfo.ID, 0, ""); err != nil {
		apiReturn.Error(c, "关闭失败")
		return
	}
	apiReturn.Success(c)
}

// 安全退出
func (l *LoginApi) Logout(c *gin.Context) {
	// userInfo, _ := base.GetCurrentUserInfo(c)
	cToken := c.GetHeader("token")
	global.CUserToken.Delete(cToken)
	apiReturn.Success(c)
}

// recordLoginLog 记录登录日志
func recordLoginLog(userId uint, username string, c *gin.Context, status int, remark string) {
	log := models.LoginLog{
		UserId:    userId,
		Username:  username,
		Ip:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    status,
		Remark:    remark,
	}
	go func() {
		if err := log.Create(); err != nil {
			global.Logger.Debug("登录日志记录失败:", err.Error())
		}
	}()
}
