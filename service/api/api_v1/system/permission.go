package system

import (
	"fmt"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type PermissionApi struct{}

func (a *PermissionApi) GetList(c *gin.Context) {
	mPermission := models.Permission{}
	list, err := mPermission.GetList()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *PermissionApi) GetByModule(c *gin.Context) {
	type Request struct {
		ModuleCode string `json:"moduleCode"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mPermission := models.Permission{}
	list, err := mPermission.GetByModule(req.ModuleCode)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, list)
}

func (a *PermissionApi) GetPermissionMatrix(c *gin.Context) {
	type Request struct {
		RoleId uint `json:"roleId"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mPermission := models.Permission{}
	permissions, err := mPermission.GetList()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	var rolePermissionIds []uint
	if req.RoleId > 0 {
		mRolePermission := models.RolePermission{}
		rolePermissionIds, err = mRolePermission.GetPermissionIdsByRoleId(req.RoleId)
		if err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	permissionIdMap := make(map[uint]bool)
	for _, pid := range rolePermissionIds {
		permissionIdMap[pid] = true
	}

	result := make([]map[string]interface{}, 0)
	for _, p := range permissions {
		result = append(result, map[string]interface{}{
			"id":             p.ID,
			"moduleCode":     p.ModuleCode,
			"moduleName":     p.ModuleName,
			"permissionId":   p.PermissionId,
			"permissionName": p.PermissionName,
			"parentId":       p.ParentId,
			"sort":           p.Sort,
			"status":         p.Status,
			"checked":        permissionIdMap[p.ID],
		})
	}

	apiReturn.SuccessData(c, result)
}

func (a *PermissionApi) SaveRolePermissions(c *gin.Context) {
	type Request struct {
		RoleId       uint   `json:"roleId"`
		PermissionIds []uint `json:"permissionIds"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.RoleId == 0 {
		apiReturn.ErrorParamFomat(c, "角色ID不能为空")
		return
	}

	mRolePermission := models.RolePermission{}
	err := mRolePermission.BatchSave(req.RoleId, req.PermissionIds)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	global.SystemSetting.Cache.Delete("role_permission_" + fmt.Sprintf("%d", req.RoleId))

	apiReturn.Success(c)
}

func (a *PermissionApi) GetRolePermissions(c *gin.Context) {
	type Request struct {
		RoleId uint `json:"roleId"`
	}
	req := Request{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.RoleId == 0 {
		apiReturn.ErrorParamFomat(c, "角色ID不能为空")
		return
	}

	mRolePermission := models.RolePermission{}
	permissionIds, err := mRolePermission.GetPermissionIdsByRoleId(req.RoleId)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, permissionIds)
}

func (a *PermissionApi) Create(c *gin.Context) {
	type Request struct {
		ModuleCode     string `json:"moduleCode" validate:"required"`
		ModuleName     string `json:"moduleName" validate:"required"`
		PermissionId   string `json:"permissionId" validate:"required"`
		PermissionName string `json:"permissionName" validate:"required"`
		ParentId       uint   `json:"parentId"`
		Sort           int    `json:"sort"`
		Status         int    `json:"status"`
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

	permission := models.Permission{
		ModuleCode:     req.ModuleCode,
		ModuleName:     req.ModuleName,
		PermissionId:   req.PermissionId,
		PermissionName: req.PermissionName,
		ParentId:       req.ParentId,
		Sort:           req.Sort,
		Status:         req.Status,
	}

	created, err := permission.Create()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, created)
}

func (a *PermissionApi) Update(c *gin.Context) {
	type Request struct {
		Id             uint   `json:"id"`
		ModuleCode     string `json:"moduleCode" validate:"required"`
		ModuleName     string `json:"moduleName" validate:"required"`
		PermissionId   string `json:"permissionId" validate:"required"`
		PermissionName string `json:"permissionName" validate:"required"`
		ParentId       uint   `json:"parentId"`
		Sort           int    `json:"sort"`
		Status         int    `json:"status"`
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

	permission := models.Permission{
		ModuleCode:     req.ModuleCode,
		ModuleName:     req.ModuleName,
		PermissionId:   req.PermissionId,
		PermissionName: req.PermissionName,
		ParentId:       req.ParentId,
		Sort:           req.Sort,
		Status:         req.Status,
	}

	err := permission.Update(req.Id)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}