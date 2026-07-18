package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type IconApi struct{}

func (a *IconApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mIcon := models.Icon{}
	icon, err := mIcon.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, icon)
}

func (a *IconApi) GetByName(c *gin.Context) {
	type Request struct {
		Name string `json:"name"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mIcon := models.Icon{}
	icon, err := mIcon.GetByName(req.Name)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, icon)
}

func (a *IconApi) GetList(c *gin.Context) {
	type Request struct {
		PageNum     int    `json:"pageNum"`
		PageSize    int    `json:"pageSize"`
		CategoryId  uint   `json:"categoryId"`
		Keyword     string `json:"keyword"`
	}
	req := Request{PageNum: 1, PageSize: 20}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mIcon := models.Icon{}
	list, total, err := mIcon.GetList(req.CategoryId, req.Keyword, req.PageNum, req.PageSize)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (a *IconApi) GetBatch(c *gin.Context) {
	type Request struct {
		Ids       []uint `json:"ids"`
		CategoryId uint  `json:"categoryId"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mIcon := models.Icon{}
	var list []models.Icon
	var err error

	if len(req.Ids) > 0 {
		list, err = mIcon.GetByIds(req.Ids)
	} else if req.CategoryId > 0 {
		list, err = mIcon.GetByCategory(req.CategoryId)
	} else {
		apiReturn.ErrorParamFomat(c, "请提供ids或categoryId")
		return
	}

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *IconApi) Create(c *gin.Context) {
	type Request struct {
		Name       string `json:"name" validate:"required"`
		IconKey    string `json:"iconKey" validate:"required"`
		CategoryId uint   `json:"categoryId"`
		Src        string `json:"src"`
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

	icon := models.Icon{
		Name:       req.Name,
		IconKey:    req.IconKey,
		CategoryId: req.CategoryId,
		Src:        req.Src,
	}

	created, err := icon.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *IconApi) Update(c *gin.Context) {
	type Request struct {
		Id         uint   `json:"id"`
		Name       string `json:"name" validate:"required"`
		IconKey    string `json:"iconKey" validate:"required"`
		CategoryId uint   `json:"categoryId"`
		Src        string `json:"src"`
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

	icon := models.Icon{
		Name:       req.Name,
		IconKey:    req.IconKey,
		CategoryId: req.CategoryId,
		Src:        req.Src,
	}

	err := icon.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *IconApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mIcon := models.Icon{}
	err := mIcon.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *IconApi) Favorite(c *gin.Context) {
	type Request struct {
		IconId uint `json:"iconId"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	mFavorite := models.IconFavorite{}
	err := mFavorite.Toggle(userInfo.ID, req.IconId)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *IconApi) GetFavorites(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	mFavorite := models.IconFavorite{}
	list, err := mFavorite.GetByUserId(userInfo.ID)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *IconApi) GetCategories(c *gin.Context) {
	mCategory := models.IconCategory{}
	list, err := mCategory.GetList()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *IconApi) CreateCategory(c *gin.Context) {
	type Request struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}
	req := Request{Sort: 0}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	category := models.IconCategory{
		Name: req.Name,
		Sort: req.Sort,
	}

	created, err := category.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *IconApi) UpdateCategory(c *gin.Context) {
	type Request struct {
		Id          uint   `json:"id"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
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

	category := models.IconCategory{
		Name: req.Name,
		Sort: req.Sort,
	}

	err := category.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *IconApi) DeleteCategory(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mCategory := models.IconCategory{}
	err := mCategory.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorByCodeAndMsg(c, 1003, err.Error())
		return
	}

	apiReturn.Success(c)
}