package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type NoticeManageApi struct{}

func (a *NoticeManageApi) GetList(c *gin.Context) {
	type Request struct {
		Page       int    `json:"page"`
		PageSize   int    `json:"pageSize"`
		Keyword    string `json:"keyword"`
		NoticeType int    `json:"noticeType"`
		Status     int    `json:"status"`
	}
	req := Request{Page: 1, PageSize: 10}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mNotice := models.Notice{}
	list, count, err := mNotice.GetList(req.Page, req.PageSize, req.Keyword, req.NoticeType, req.Status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

func (a *NoticeManageApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mNotice := models.Notice{}
	notice, err := mNotice.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, notice)
}

func (a *NoticeManageApi) Create(c *gin.Context) {
	type Request struct {
		Title         string `json:"title" validate:"required,min=1,max=255"`
		Content       string `json:"content" validate:"required"`
		DisplayType   int    `json:"displayType"`
		NoticeType    int    `json:"noticeType"`
		OneRead       int    `json:"oneRead"`
		Url           string `json:"url"`
		IsLogin       uint   `json:"isLogin"`
		Status        int    `json:"status"`
		TargetUserIds string `json:"targetUserIds"`
	}
	req := Request{DisplayType: 2, NoticeType: 1, Status: 1}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	notice := models.Notice{
		Title:         req.Title,
		Content:       req.Content,
		DisplayType:   req.DisplayType,
		NoticeType:    req.NoticeType,
		OneRead:       req.OneRead,
		Url:           req.Url,
		IsLogin:       req.IsLogin,
		Status:        req.Status,
		UserId:        userInfo.ID,
		TargetUserIds: req.TargetUserIds,
	}

	created, err := notice.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *NoticeManageApi) Update(c *gin.Context) {
	type Request struct {
		Id            uint   `json:"id" validate:"required"`
		Title         string `json:"title" validate:"required,min=1,max=255"`
		Content       string `json:"content" validate:"required"`
		DisplayType   int    `json:"displayType"`
		NoticeType    int    `json:"noticeType"`
		OneRead       int    `json:"oneRead"`
		Url           string `json:"url"`
		IsLogin       uint   `json:"isLogin"`
		Status        int    `json:"status"`
		TargetUserIds string `json:"targetUserIds"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	notice := models.Notice{
		Title:         req.Title,
		Content:       req.Content,
		DisplayType:   req.DisplayType,
		NoticeType:    req.NoticeType,
		OneRead:       req.OneRead,
		Url:           req.Url,
		IsLogin:       req.IsLogin,
		Status:        req.Status,
		TargetUserIds: req.TargetUserIds,
	}

	if err := notice.Update(req.Id); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *NoticeManageApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mNotice := models.Notice{}
	if err := mNotice.Delete(req.Id); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// GetVisibleNotices 获取当前用户可见的通知（公告弹窗用）
func (a *NoticeManageApi) GetVisibleNotices(c *gin.Context) {
	type Request struct {
		DisplayType []int `json:"displayType"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	// 使用公开用户的 ID（如果是公开模式）
	userId := userInfo.ID

	mNotice := models.Notice{}
	list, err := mNotice.GetVisibleNotices(userId, req.DisplayType...)
	if err != nil {
		global.Db.Find(&list, "display_type IN ? AND status=1", req.DisplayType)
	}

	apiReturn.SuccessListData(c, list, 0)
}

// MarkRead 标记当前用户已读某公告（服务端持久化，跨设备同步）
func (a *NoticeManageApi) MarkRead(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	if req.Id == 0 {
		apiReturn.ErrorParamFomat(c, "id is required")
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mRead := models.NoticeRead{}
	if err := mRead.MarkRead(userInfo.ID, req.Id); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}