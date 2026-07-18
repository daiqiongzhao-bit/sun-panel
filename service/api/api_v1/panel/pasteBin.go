package panel

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
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

type PasteBinApi struct{}

// 生成6位随机访问码
func generatePasteCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code += string(chars[n.Int64()])
	}
	return code
}

func (a *PasteBinApi) Create(c *gin.Context) {
	type Request struct {
		Title         string `json:"title"`
		Type          int    `json:"type"`     // 1文本 2文件
		Content       string `json:"content"`  // 文本内容
		FileName      string `json:"fileName"` // 文件名
		FileSize      int64  `json:"fileSize"`
		ExpireH       int    `json:"expireH"`       // 过期小时数，默认24
		Password      string `json:"password"`      // 访问密码，空=无密码
		BurnAfterRead int    `json:"burnAfterRead"` // 1阅后即焚
	}
	req := Request{Type: 1, ExpireH: 24}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	expireAt := time.Now().Add(time.Duration(req.ExpireH) * time.Hour)

	// 密码 SHA256 哈希
	passwordHash := ""
	if req.Password != "" {
		hash := sha256.Sum256([]byte(req.Password))
		passwordHash = fmt.Sprintf("%x", hash)
	}

	item := models.PasteBin{
		UserId:        userInfo.ID,
		Type:          req.Type,
		Title:         req.Title,
		Content:       req.Content,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		Code:          generatePasteCode(),
		Password:      passwordHash,
		BurnAfterRead: req.BurnAfterRead,
		ExpireAt:      expireAt,
		Status:        1,
	}

	created, err := item.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessData(c, created)
}

func (a *PasteBinApi) GetByCode(c *gin.Context) {
	type Request struct {
		Code     string `json:"code" validate:"required"`
		Password string `json:"password"` // 访问密码
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	item, err := (&models.PasteBin{}).GetByCode(req.Code)
	if err != nil {
		apiReturn.ErrorCode(c, 1003, "内容不存在或已过期", nil)
		return
	}

	// 密码验证
	if item.Password != "" {
		hash := sha256.Sum256([]byte(req.Password))
		providedHash := fmt.Sprintf("%x", hash)
		if providedHash != item.Password {
			apiReturn.ErrorCode(c, 1004, "访问密码错误", nil)
			return
		}
	}

	// 增加访问计数
	(&models.PasteBin{}).IncrementAccess(item.ID)

	// 阅后即焚：读取后立即物理删除
	if item.BurnAfterRead == 1 {
		global.Db.Delete(&item)
	}

	// 不返回 userId
	item.UserId = 0
	apiReturn.SuccessData(c, item)
}

func (a *PasteBinApi) GetMyList(c *gin.Context) {
	type Request struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	list, count, err := (&models.PasteBin{}).GetByUser(userInfo.ID, req.Page, req.PageSize)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessListData(c, list, count)
}

func (a *PasteBinApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id" validate:"required"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := (&models.PasteBin{}).Delete(req.Id, userInfo.ID); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// Update 修改中转记录（密码、过期时间等）
func (a *PasteBinApi) Update(c *gin.Context) {
	type Request struct {
		Id            uint   `json:"id" validate:"required"`
		Password      string `json:"password"`
		ExpireH       int    `json:"expireH"`
		BurnAfterRead int    `json:"burnAfterRead"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	updates := make(map[string]interface{})

	if req.Password != "" {
		hash := sha256.Sum256([]byte(req.Password))
		updates["password"] = fmt.Sprintf("%x", hash)
	} else {
		updates["password"] = ""
	}

	if req.ExpireH > 0 {
		expireAt := time.Now().Add(time.Duration(req.ExpireH) * time.Hour)
		updates["expire_at"] = expireAt
	}

	updates["burn_after_read"] = req.BurnAfterRead

	if err := (&models.PasteBin{}).Update(req.Id, userInfo.ID, updates); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// UploadFile 上传中转站文件（支持任意文件类型）
func (a *PasteBinApi) UploadFile(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	fileExt := strings.ToLower(path.Ext(f.Filename))
	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fildDir := fmt.Sprintf("%s/pasteBin/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(fildDir)
	if !isExist {
		os.MkdirAll(fildDir, os.ModePerm)
	}
	filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
	if err := c.SaveUploadedFile(f, filepath); err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}

	// 记录到数据库
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, f.Filename, fileExt, filepath)
	apiReturn.SuccessData(c, gin.H{
		"url": filepath[1:], // 去掉前导斜杠
	})
}

// DownloadFile 下载中转站文件（支持 GET 和 POST）
func (a *PasteBinApi) DownloadFile(c *gin.Context) {
	code := c.PostForm("code")
	if code == "" {
		code = c.Query("code")
	}
	if code == "" {
		// try JSON body
		type Request struct {
			Code     string `json:"code"`
			Password string `json:"password"`
		}
		req := Request{}
		_ = c.ShouldBindBodyWith(&req, binding.JSON)
		code = req.Code
	}
	if code == "" {
		apiReturn.ErrorParamFomat(c, "code is required")
		return
	}

	item, err := (&models.PasteBin{}).GetByCode(code)
	if err != nil || item.Type != 2 {
		apiReturn.ErrorCode(c, 1202, "文件不存在或已过期", nil)
		return
	}

	// 密码验证
	if item.Password != "" {
		password := c.PostForm("password")
		if password == "" {
			password = c.Query("password")
		}
		hash := sha256.Sum256([]byte(password))
		if fmt.Sprintf("%x", hash) != item.Password {
			apiReturn.ErrorCode(c, 1004, "访问密码错误", nil)
			return
		}
	}

	configUpload := global.Config.GetValueString("base", "source_path")
	filePath := configUpload + "/" + item.Content
	if exists, _ := cmn.PathExists(filePath); !exists {
		apiReturn.ErrorCode(c, 1202, "文件不存在", nil)
		return
	}

	// 增加计数
	(&models.PasteBin{}).IncrementAccess(item.ID)

	c.File(filePath)
}

// GetAccessUrl 获取访问URL（供前端生成链接/二维码）
func (a *PasteBinApi) GetAccessUrl(c *gin.Context) {
	siteUrl, _ := global.SystemSetting.GetValueString("system")
	if siteUrl == "" {
		siteUrl = "http://localhost:3030"
	}
	type Request struct {
		Code string `json:"code" validate:"required"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	url := fmt.Sprintf("%s/#/paste/%s", siteUrl, req.Code)
	apiReturn.SuccessData(c, gin.H{"url": url, "code": req.Code})
}