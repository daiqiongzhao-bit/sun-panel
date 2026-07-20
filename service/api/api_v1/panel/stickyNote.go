package panel

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type StickyNoteApi struct{}

func (a *StickyNoteApi) GetList(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	notes, err := (&models.StickyNote{}).GetByUser(userInfo.ID)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessListData(c, notes, 0)
}

func (a *StickyNoteApi) Create(c *gin.Context) {
	req := models.StickyNote{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	req.UserId = userInfo.ID
	if req.Color == "" {
		req.Color = "#fff3bf"
	}
	if req.Width == 0 {
		req.Width = 200
	}
	if req.Height == 0 {
		req.Height = 150
	}

	created, err := req.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessData(c, created)
}

func (a *StickyNoteApi) Update(c *gin.Context) {
	type Request struct {
		Id      uint   `json:"id" validate:"required"`
		Content string `json:"content"`
		Color   string `json:"color"`
		PosX    int    `json:"posX"`
		PosY    int    `json:"posY"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		ZIndex  int    `json:"zIndex"`
		Status  *int   `json:"status"` // 1启用 0停用
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	updates := map[string]interface{}{
		"content": req.Content,
		"color":   req.Color,
		"pos_x":   req.PosX,
		"pos_y":   req.PosY,
		"width":   req.Width,
		"height":  req.Height,
		"z_index": req.ZIndex,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	note := models.StickyNote{
		Content: req.Content,
		Color:   req.Color,
		PosX:    req.PosX,
		PosY:    req.PosY,
		Width:   req.Width,
		Height:  req.Height,
		ZIndex:  req.ZIndex,
	}
	if err := note.UpdateWithMap(req.Id, userInfo.ID, updates); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a *StickyNoteApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id" validate:"required"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	if err := (&models.StickyNote{}).Delete(req.Id, userInfo.ID); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}