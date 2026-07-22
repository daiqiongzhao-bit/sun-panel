package panel

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/panelApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/department"
	"sun-panel/lib/siteFavicon"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type ItemIcon struct {
}

func (a *ItemIcon) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ItemIconGroupId == 0 {
		// apiReturn.Error(c, "Group is mandatory")
		apiReturn.ErrorParamFomat(c, "Group is mandatory")
		return
	}

	// json转字符串
	if j, err := json.Marshal(req.Icon); err == nil {
		req.IconJson = string(j)
	}

	if req.ID != 0 {
		// 修改：先校验归属，非管理员只能改自己的图标
		existing := models.ItemIcon{}
		if err := global.Db.First(&existing, "id=?", req.ID).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		if userInfo.Role != department.ROLE_SUPER_ADMIN && existing.UserId != userInfo.ID {
			apiReturn.ErrorNoAccess(c)
			return
		}
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId", "ItemIconGroupId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&models.ItemIcon{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		req.Sort = 9999
		// 创建：归属当前用户，自动继承部门
		req.UserId = userInfo.ID
		req.DepartmentId = userInfo.DepartmentId
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

// 添加多个图标
func (a *ItemIcon) AddMultiple(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	// type Request
	req := []models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	for i := 0; i < len(req); i++ {
		if req[i].ItemIconGroupId == 0 {
			apiReturn.ErrorParamFomat(c, "Group is mandatory")
			return
		}
		req[i].UserId = userInfo.ID
		req[i].DepartmentId = userInfo.DepartmentId // 自动继承用户的部门
		// json转字符串
		if j, err := json.Marshal(req[i].Icon); err == nil {
			req[i].IconJson = string(j)
		}
	}

	global.Db.Create(&req)

	apiReturn.SuccessData(c, req)
}

// // 获取详情
// func (a *ItemIcon) GetInfo(c *gin.Context) {
// 	req := systemApiStructs.AiDrawGetInfoReq{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	userInfo, _ := base.GetCurrentUserInfo(c)

// 	aiDraw := models.AiDraw{}
// 	aiDraw.ID = req.ID
// 	if err := aiDraw.GetInfo(global.Db); err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			apiReturn.Error(c, "不存在记录")
// 			return
// 		}
// 		apiReturn.ErrorDatabase(c, err.Error())
// 		return
// 	}

// 	if userInfo.ID != aiDraw.UserID {
// 		apiReturn.ErrorNoAccess(c)
// 		return
// 	}

// 	apiReturn.SuccessData(c, aiDraw)
// }

func (a *ItemIcon) GetListByGroupId(c *gin.Context) {
	req := models.ItemIcon{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	visitMode := base.GetCurrentVisitMode(c)
	itemIcons := []models.ItemIcon{}

	query := global.Db.Order("sort ,created_at")

	if visitMode == base.VISIT_MODE_PUBLIC {
		// 公开访问模式：保持原逻辑，只查公开用户的个人数据
		query = query.Where("item_icon_group_id = ? AND user_id=?", req.ItemIconGroupId, userInfo.ID)
	} else {
		// 登录模式：应用部门数据隔离
		query = query.Where("item_icon_group_id = ?", req.ItemIconGroupId)
		query = department.BuildDepartmentScope(query, userInfo.ID, userInfo)
	}

	if err := query.Find(&itemIcons).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for k, v := range itemIcons {
		json.Unmarshal([]byte(v.IconJson), &itemIcons[k].Icon)
	}

	apiReturn.SuccessListData(c, itemIcons, 0)
}

func (a *ItemIcon) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	// 超级管理员可以删所有；其他用户只能删自己创建的
	q := global.Db.Where("id in ?", req.Ids)
	if userInfo.Role != department.ROLE_SUPER_ADMIN {
		q = q.Where("user_id=?", userInfo.ID)
	}
	deleteQuery := q.Delete(&models.ItemIcon{})
	if deleteQuery.Error != nil {
		apiReturn.ErrorDatabase(c, deleteQuery.Error.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIcon) SaveSort(c *gin.Context) {
	req := panelApiStructs.ItemIconSaveSortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	// 超级管理员可排序任意图标；其他用户只能排序自己的
	scopeUserId := userInfo.ID
	if userInfo.Role == department.ROLE_SUPER_ADMIN {
		scopeUserId = 0
	}

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		for _, v := range req.SortItems {
			q := tx.Model(&models.ItemIcon{}).Where("id=? AND item_icon_group_id=?", v.Id, req.ItemIconGroupId)
			if scopeUserId > 0 {
				q = q.Where("user_id=?", scopeUserId)
			}
			if err := q.Update("sort", v.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}

// 支持获取并直接下载对方网站图标到服务器
func (a *ItemIcon) GetSiteFavicon(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := panelApiStructs.ItemIconGetSiteFaviconReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	resp := panelApiStructs.ItemIconGetSiteFaviconResp{}
	// 本地缓存：同一主机已获取过且文件仍在，则直接复用，不再访问外网（减少依赖、加速）
	if hu, herr := url.Parse(req.Url); herr == nil && hu.Host != "" {
		var cached models.File
		if global.Db.Where("category=? AND file_name=? AND user_id=?", "icon", hu.Host, userInfo.ID).First(&cached).Error == nil && cached.Src != "" {
			cfg := global.Config.GetValueString("base", "source_path")
			if exists, _ := cmn.PathExists(cfg + cached.Src); exists {
				resp.IconUrl = cached.Src
				apiReturn.SuccessData(c, resp)
				return
			}
		}
	}
	fullUrl, favErr := siteFavicon.GetOneFaviconURL(req.Url)
	if favErr != nil {
		apiReturn.Error(c, "acquisition failed: get ico error:"+favErr.Error())
		return
	}
	if fullUrl == "" {
		// 所有公共 favicon 服务均不可达：返回空图标地址，前端据此显示占位图标，不再报错
		resp.IconUrl = ""
		apiReturn.SuccessData(c, resp)
		return
	}

	parsedURL, err := url.Parse(req.Url)
	if err != nil {
		apiReturn.Error(c, "acquisition failed:"+err.Error())
		return
	}

	protocol := parsedURL.Scheme
	global.Logger.Debug("protocol:", protocol)
	global.Logger.Debug("fullUrl:", fullUrl)

	// 如果URL以双斜杠（//）开头，则使用当前页面协议
	if strings.HasPrefix(fullUrl, "//") {
		fullUrl = protocol + "://" + fullUrl[2:]
	} else if !strings.HasPrefix(fullUrl, "http://") && !strings.HasPrefix(fullUrl, "https://") {
		// 如果URL既不以http://开头也不以https://开头，则默认为http协议
		fullUrl = "http://" + fullUrl
	}
	global.Logger.Debug("fullUrl:", fullUrl)
	// 去除图标的get参数
	{
		parsedIcoURL, err := url.Parse(fullUrl)
		if err != nil {
			apiReturn.Error(c, "acquisition failed: parsed ico URL :"+err.Error())
			return
		}
		fullUrl = parsedIcoURL.Scheme + "://" + parsedIcoURL.Host + parsedIcoURL.Path
	}
	global.Logger.Debug("fullUrl:", fullUrl)

	// 生成保存目录
	configUpload := global.Config.GetValueString("base", "source_path")
	savePath := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(savePath)
	if !isExist {
		os.MkdirAll(savePath, os.ModePerm)
	}

	// 下载
	var imgInfo *os.File
	{
		var err error
		if imgInfo, err = siteFavicon.DownloadImage(fullUrl, savePath, 1024*1024); err != nil {
			apiReturn.Error(c, "acquisition failed: download"+err.Error())
			return
		}
	}

	// 保存到数据库（扩展名去掉查询串/锚点，无法识别时默认 .ico，确保 ico/svg 被正确归类）
	cleanFullURL := fullUrl
	if idx := strings.IndexAny(cleanFullURL, "?#"); idx >= 0 {
		cleanFullURL = cleanFullURL[:idx]
	}
	ext := strings.ToLower(path.Ext(cleanFullURL))
	if ext == "" {
		ext = ".ico"
	}
	mFile := models.File{}
	if _, err := mFile.AddFile(userInfo.ID, parsedURL.Host, ext, imgInfo.Name(), "icon"); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	resp.IconUrl = imgInfo.Name()[1:]
	apiReturn.SuccessData(c, resp)
}
