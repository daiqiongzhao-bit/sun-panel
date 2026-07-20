package panel

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SearchApi struct{}

// LocalSearch 本地资产模糊搜索（图标、分组、中转站）
func (a *SearchApi) LocalSearch(c *gin.Context) {
	type Request struct {
		Keyword string `json:"keyword"`
		Limit   int    `json:"limit"`
	}
	req := Request{Limit: 10}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Keyword == "" {
		apiReturn.SuccessListData(c, []interface{}{}, 0)
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	keyword := "%" + req.Keyword + "%"

	type Result struct {
		Type  string `json:"type"`  // icon/group/pastebin
		Title string `json:"title"`
		Url   string `json:"url"`
		Desc  string `json:"desc"`
		Code  string `json:"code"` // pastebin专用
	}

	results := []Result{}

	// 1. 搜索图标（按用户可见范围）
	var icons []models.ItemIcon
	global.Db.Where("(title LIKE ? OR description LIKE ? OR url LIKE ?) AND user_id=?",
		keyword, keyword, keyword, userInfo.ID).
		Limit(req.Limit).Find(&icons)
	for _, icon := range icons {
		results = append(results, Result{
			Type:  "icon",
			Title: icon.Title,
			Url:   icon.Url,
			Desc:  icon.Description,
		})
	}

	// 2. 搜索图标分组
	var groups []models.ItemIconGroup
	global.Db.Where("title LIKE ? AND user_id=?", keyword, userInfo.ID).
		Limit(req.Limit).Find(&groups)
	for _, g := range groups {
		results = append(results, Result{
			Type:  "group",
			Title: g.Title,
			Desc:  g.Description,
		})
	}

	// 3. 搜索中转站（仅自己的）
	if len(results) < req.Limit*2 {
		var pastes []models.PasteBin
		global.Db.Where("(title LIKE ? OR content LIKE ?) AND user_id=? AND status=1",
			keyword, keyword, userInfo.ID).
			Limit(req.Limit).Find(&pastes)
		for _, p := range pastes {
			results = append(results, Result{
				Type:  "pastebin",
				Title: p.Title,
				Desc:  truncate(p.Content, 100),
				Code:  p.Code,
			})
		}
	}

	// 限制总数
	if len(results) > req.Limit*2 {
		results = results[:req.Limit*2]
	}

	apiReturn.SuccessListData(c, results, int64(len(results)))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}