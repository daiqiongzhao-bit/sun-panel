package openness

import (
	"strings"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
)

type Openness struct {
}

func (a *Openness) LoginConfig(c *gin.Context) {
	cfg := systemSetting.ApplicationSetting{}
	if err := global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &cfg); err != nil {
		apiReturn.Error(c, "配置查询失败："+err.Error())
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"loginCaptcha": cfg.LoginCaptcha,
		"loginAllowIps": cfg.Login.LoginAllowIps,
		"register":     cfg.Register,
	})
}

// SetLoginConfig 保存登录安全相关配置（验证码开关、允许登录的 IP 白名单）
func (a *Openness) SetLoginConfig(c *gin.Context) {
	var req struct {
		LoginCaptcha bool   `json:"loginCaptcha"`
		LoginAllowIps string `json:"loginAllowIps"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.Error(c, "参数错误："+err.Error())
		return
	}

	cfg := systemSetting.ApplicationSetting{}
	// 保留已有的其它字段（如 register、webSiteUrl）
	_ = global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &cfg)
	cfg.LoginCaptcha = req.LoginCaptcha
	cfg.Login.LoginAllowIps = strings.TrimSpace(req.LoginAllowIps)

	if err := global.SystemSetting.Set(systemSetting.SYSTEM_APPLICATION, cfg); err != nil {
		apiReturn.Error(c, "保存失败："+err.Error())
		return
	}
	apiReturn.SuccessData(c, "保存成功")
}

func (a *Openness) GetDisclaimer(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.DISCLAIMER); err != nil {
		global.SystemSetting.Set(systemSetting.DISCLAIMER, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}

func (a *Openness) GetAboutDescription(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.WEB_ABOUT_DESCRIPTION); err != nil {
		global.SystemSetting.Set(systemSetting.WEB_ABOUT_DESCRIPTION, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}
