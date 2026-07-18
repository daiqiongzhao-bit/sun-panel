package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LogApi struct{}

// 操作日志列表
func (a *LogApi) GetOperationLogList(c *gin.Context) {
	type Request struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Keyword  string `json:"keyword"`
		Module   string `json:"module"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mLog := models.OperationLog{}
	list, count, err := mLog.GetList(req.Page, req.PageSize, req.Keyword, req.Module)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

// 登录日志列表
func (a *LogApi) GetLoginLogList(c *gin.Context) {
	type Request struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Keyword  string `json:"keyword"`
	}
	req := Request{Page: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mLog := models.LoginLog{}
	list, count, err := mLog.GetList(req.Page, req.PageSize, req.Keyword)
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