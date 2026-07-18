package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RoleApi struct{}

func (a *RoleApi) GetList(c *gin.Context) {
	type Request struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Keyword  string `json:"keyword"`
	}
	req := Request{
		Page:     1,
		PageSize: 10,
	}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mRole := models.Role{}
	var list []models.Role
	var count int64
	db := models.Db.Model(mRole)
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	err := db.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Order("id asc").Find(&list).Error
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, count)
}

func (a *RoleApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mRole := models.Role{}
	role, err := mRole.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, role)
}

func (a *RoleApi) Create(c *gin.Context) {
	type Request struct {
		Name        string `json:"name" validate:"required,min=2,max=50"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}
	req := Request{Status: 1}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if errMsg, err := base.ValidateInputStruct(&req); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mRole := models.Role{}
	if mRole.CheckNameExists(req.Name, 0) {
		apiReturn.ErrorByCodeAndMsg(c, 1001, "角色名称已存在")
		return
	}

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	created, err := role.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *RoleApi) Update(c *gin.Context) {
	type Request struct {
		Id          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// Name 和 Description 字段校验（仅当提供了这两个字段时才校验）
	if req.Name != "" && (len(req.Name) < 2 || len(req.Name) > 50) {
		apiReturn.ErrorParamFomat(c, "name长度应在2-50个字符之间")
		return
	}

	mRole := models.Role{}
	if req.Name != "" && mRole.CheckNameExists(req.Name, req.Id) {
		apiReturn.ErrorByCodeAndMsg(c, 1001, "角色名称已存在")
		return
	}

	updateMap := map[string]interface{}{}
	if req.Name != "" {
		updateMap["name"] = req.Name
	}
	if req.Description != "" {
		updateMap["description"] = req.Description
	}
	if req.Status != nil {
		updateMap["status"] = *req.Status
	}

	if len(updateMap) == 0 {
		apiReturn.ErrorParamFomat(c, "没有需要更新的字段")
		return
	}

	err := mRole.UpdateFields(req.Id, updateMap)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *RoleApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mRole := models.Role{}
	err := mRole.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorByCodeAndMsg(c, 1002, err.Error())
		return
	}

	apiReturn.Success(c)
}