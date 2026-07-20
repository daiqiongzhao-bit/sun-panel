package system

import (
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type NoticeApi struct {
}

func (a *NoticeApi) GetListByDisplayType(c *gin.Context) {
	req := systemApiStructs.NoticeGetListByDisplayTypeReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	noticeList := []models.Notice{}
	// 仅返回启用公告，且排除当前用户已读(oneRead=1)的公告
	if err := global.Db.Where("display_type IN ? AND status=? AND (one_read=0 OR id NOT IN (SELECT notice_id FROM notice_read WHERE user_id=?))",
		req.DisplayType, models.NOTICE_STATUS_ENABLED, userInfo.ID).Order("id desc").Find(&noticeList).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessListData(c, noticeList, 0)
}
