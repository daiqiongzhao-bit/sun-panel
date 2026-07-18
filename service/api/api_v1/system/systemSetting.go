package system

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemSettingApi struct{}

const (
	SYS_SETTING_LOGO           = "system_logo"
	SYS_SETTING_LOGIN_BACKGROUND = "login_background"
)

type LogoConfig struct {
	ImageUrl   string `json:"imageUrl"`
	Size       int    `json:"size"`
	UseCDN     bool   `json:"useCDN"`
	CDNUrl     string `json:"cdnUrl"`
}

type BackgroundConfig struct {
	ImageUrl      string `json:"imageUrl"`
	DisplayMode   string `json:"displayMode"`
	UseCustomUrl  bool   `json:"useCustomUrl"`
	CustomUrl     string `json:"customUrl"`
}

func (a *SystemSettingApi) GetLogoConfig(c *gin.Context) {
	var config LogoConfig
	err := global.SystemSetting.GetValueByInterface(SYS_SETTING_LOGO, &config)
	if err != nil {
		config = LogoConfig{
			ImageUrl:   "/assets/logo.png",
			Size:       80,
			UseCDN:     false,
			CDNUrl:     "",
		}
	}
	apiReturn.SuccessData(c, config)
}

func (a *SystemSettingApi) SetLogoConfig(c *gin.Context) {
	type Request struct {
		ImageUrl   string `json:"imageUrl"`
		Size       int    `json:"size"`
		UseCDN     bool   `json:"useCDN"`
		CDNUrl     string `json:"cdnUrl"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Size < 50 || req.Size > 200 {
		apiReturn.ErrorParamFomat(c, "Logo尺寸范围为50-200px")
		return
	}

	config := LogoConfig{
		ImageUrl: req.ImageUrl,
		Size:     req.Size,
		UseCDN:   req.UseCDN,
		CDNUrl:   req.CDNUrl,
	}

	err := global.SystemSetting.Set(SYS_SETTING_LOGO, config)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *SystemSettingApi) UploadLogo(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	agreeExts := []string{".png", ".jpg", ".jpeg"}
	if !cmn.InArray(agreeExts, fileExt) {
		apiReturn.ErrorByCode(c, 1301)
		return
	}

	fileSize := f.Size
	if fileSize > 5*1024*1024 {
		apiReturn.ErrorByCodeAndMsg(c, 1302, "文件大小不能超过5MB")
		return
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/logo/", configUpload)
	isExist, _ := cmn.PathExists(fileDir)
	if !isExist {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	filePath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	c.SaveUploadedFile(f, filePath)

	mFile := models.File{}
	mFile.AddFile(userInfo.ID, f.Filename, fileExt, filePath)

	apiReturn.SuccessData(c, gin.H{
		"imageUrl": filePath[1:],
	})
}

func (a *SystemSettingApi) GetBackgroundConfig(c *gin.Context) {
	var config BackgroundConfig
	err := global.SystemSetting.GetValueByInterface(SYS_SETTING_LOGIN_BACKGROUND, &config)
	if err != nil {
		config = BackgroundConfig{
			ImageUrl:      "/assets/defaultBackground.webp",
			DisplayMode:   "cover",
			UseCustomUrl:  false,
			CustomUrl:     "",
		}
	}
	apiReturn.SuccessData(c, config)
}

func (a *SystemSettingApi) SetBackgroundConfig(c *gin.Context) {
	type Request struct {
		ImageUrl      string `json:"imageUrl"`
		DisplayMode   string `json:"displayMode"`
		UseCustomUrl  bool   `json:"useCustomUrl"`
		CustomUrl     string `json:"customUrl"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	validModes := []string{"cover", "repeat", "center", "stretch"}
	if !cmn.InArray(validModes, req.DisplayMode) {
		apiReturn.ErrorParamFomat(c, "显示方式只能是cover/repeat/center/stretch")
		return
	}

	config := BackgroundConfig{
		ImageUrl:     req.ImageUrl,
		DisplayMode:  req.DisplayMode,
		UseCustomUrl: req.UseCustomUrl,
		CustomUrl:    req.CustomUrl,
	}

	err := global.SystemSetting.Set(SYS_SETTING_LOGIN_BACKGROUND, config)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *SystemSettingApi) UploadBackground(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	agreeExts := []string{".png", ".jpg", ".jpeg"}
	if !cmn.InArray(agreeExts, fileExt) {
		apiReturn.ErrorByCode(c, 1301)
		return
	}

	fileSize := f.Size
	if fileSize > 5*1024*1024 {
		apiReturn.ErrorByCodeAndMsg(c, 1302, "文件大小不能超过5MB")
		return
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/background/", configUpload)
	isExist, _ := cmn.PathExists(fileDir)
	if !isExist {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	filePath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	c.SaveUploadedFile(f, filePath)

	mFile := models.File{}
	mFile.AddFile(userInfo.ID, f.Filename, fileExt, filePath)

	apiReturn.SuccessData(c, gin.H{
		"imageUrl": filePath[1:],
	})
}

func (a *SystemSettingApi) GetPresetBackgrounds(c *gin.Context) {
	presets := []map[string]interface{}{
		{
			"id":    1,
			"name":  "星空",
			"imageUrl": "/assets/defaultBackground.webp",
		},
		{
			"id":    2,
			"name":  "蓝天",
			"imageUrl": "/assets/start_sky.jpg",
		},
		{
			"id":    3,
			"name":  "渐变紫",
			"imageUrl": "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
			"isGradient": true,
		},
		{
			"id":    4,
			"name":  "深海",
			"imageUrl": "linear-gradient(135deg, #0c1445 0%, #1a237e 50%, #0d47a1 100%)",
			"isGradient": true,
		},
		{
			"id":    5,
			"name":  "暖阳",
			"imageUrl": "linear-gradient(135deg, #ff9a9e 0%, #fad0c4 99%, #fad0c4 100%)",
			"isGradient": true,
		},
		{
			"id":    6,
			"name":  "森林",
			"imageUrl": "linear-gradient(135deg, #134e5e 0%, #71b280 100%)",
			"isGradient": true,
		},
	}
	apiReturn.SuccessData(c, presets)
}