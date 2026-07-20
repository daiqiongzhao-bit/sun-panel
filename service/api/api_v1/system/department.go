package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/department"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DepartmentApi struct{}

func (a *DepartmentApi) GetList(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	mDepartment := models.Department{}
	query := global.Db.Model(&mDepartment)

	// 数据隔离：非超级管理员仅能查看可见部门（部门管理员=本部+子部门，普通用户=本部）
	if userInfo.Role != department.ROLE_SUPER_ADMIN {
		if deptIds, _ := department.GetUserVisibleDepartmentIds(userInfo); deptIds != nil {
			query = query.Where("id IN ?", deptIds)
		}
	}

	list := []models.Department{}
	if err := query.Find(&list).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *DepartmentApi) GetTreeList(c *gin.Context) {
	mDepartment := models.Department{}
	list, err := mDepartment.GetTreeList()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	tree := buildDepartmentTree(list, 0)
	apiReturn.SuccessData(c, tree)
}

func buildDepartmentTree(list []models.Department, parentId uint) []map[string]interface{} {
	tree := make([]map[string]interface{}, 0)
	for _, dept := range list {
		if dept.ParentId == parentId {
			node := map[string]interface{}{
				"id":          dept.ID,
				"name":        dept.Name,
				"code":        dept.Code,
				"parentId":    dept.ParentId,
				"description": dept.Description,
				"leaderId":    dept.LeaderId,
				"status":      dept.Status,
				"sort":        dept.Sort,
				"children":    buildDepartmentTree(list, dept.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

func (a *DepartmentApi) GetById(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mDepartment := models.Department{}
	dept, err := mDepartment.GetById(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, dept)
}

func (a *DepartmentApi) Create(c *gin.Context) {
	type Request struct {
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Code        string `json:"code" validate:"required,min=2,max=50"`
		ParentId    uint   `json:"parentId"`
		Description string `json:"description"`
		LeaderId    uint   `json:"leaderId"`
		Status      int    `json:"status"`
		Sort        int    `json:"sort"`
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

	department := models.Department{
		Name:        req.Name,
		Code:        req.Code,
		ParentId:    req.ParentId,
		Description: req.Description,
		LeaderId:    req.LeaderId,
		Status:      req.Status,
		Sort:        req.Sort,
	}

	created, err := department.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *DepartmentApi) Update(c *gin.Context) {
	type Request struct {
		Id          uint   `json:"id"`
		Name        string `json:"name" validate:"required,min=2,max=100"`
		Code        string `json:"code" validate:"required,min=2,max=50"`
		ParentId    uint   `json:"parentId"`
		Description string `json:"description"`
		LeaderId    uint   `json:"leaderId"`
		Status      int    `json:"status"`
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

	department := models.Department{
		Name:        req.Name,
		Code:        req.Code,
		ParentId:    req.ParentId,
		Description: req.Description,
		LeaderId:    req.LeaderId,
		Status:      req.Status,
		Sort:        req.Sort,
	}

	err := department.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a *DepartmentApi) Delete(c *gin.Context) {
	type Request struct {
		Id uint `json:"id"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mDepartment := models.Department{}
	err := mDepartment.Delete(req.Id)
	if err != nil {
		apiReturn.ErrorByCodeAndMsg(c, 1003, err.Error())
		return
	}

	apiReturn.Success(c)
}
