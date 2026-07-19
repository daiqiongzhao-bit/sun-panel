package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LogApi struct{}

// 操作日志列表
func (a *LogApi) GetOperationLogList(c *gin.Context) {
	type Request struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		Keyword   string `json:"keyword"`
		Module    string `json:"module"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var startTime, endTime *time.Time
	if req.StartTime > 0 {
		t := time.Unix(req.StartTime, 0)
		startTime = &t
	}
	if req.EndTime > 0 {
		t := time.Unix(req.EndTime, 0)
		endTime = &t
	}

	mLog := models.OperationLog{}
	list, count, err := mLog.GetList(req.Page, req.PageSize, req.Keyword, req.Module, startTime, endTime)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

// 登录日志列表
func (a *LogApi) GetLoginLogList(c *gin.Context) {
	type Request struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		Keyword   string `json:"keyword"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var startTime, endTime *time.Time
	if req.StartTime > 0 {
		t := time.Unix(req.StartTime, 0)
		startTime = &t
	}
	if req.EndTime > 0 {
		t := time.Unix(req.EndTime, 0)
		endTime = &t
	}

	mLog := models.LoginLog{}
	list, count, err := mLog.GetList(req.Page, req.PageSize, req.Keyword, startTime, endTime)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

// 清除操作日志
func (a *LogApi) ClearOperationLog(c *gin.Context) {
	mLog := models.OperationLog{}
	if err := mLog.ClearAll(); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// 清除登录日志
func (a *LogApi) ClearLoginLog(c *gin.Context) {
	mLog := models.LoginLog{}
	if err := mLog.ClearAll(); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}