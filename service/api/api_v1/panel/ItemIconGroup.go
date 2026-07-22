package panel

import (
	"errors"
	"math"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/department"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// errGroupNotOwned 表示当前用户不拥有待删除的分组（或分组不存在），
// 例如试图删除共享分组（user_id 为 NULL）的非管理员用户。用于 Deletes 明确返回无权限错误，
// 避免“删除接口返回成功、但分组刷新后又出现”的误导。
var errGroupNotOwned = errors.New("the requested group(s) do not belong to the current user or do not exist")

type ItemIconGroup struct {
}

func (a *ItemIconGroup) Edit(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	req := models.ItemIconGroup{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ID != 0 {
		// 修改：先校验归属，非管理员只能改自己的分组
		existing := models.ItemIconGroup{}
		if err := global.Db.First(&existing, "id=?", req.ID).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		if userInfo.Role != department.ROLE_SUPER_ADMIN && existing.UserId != userInfo.ID {
			apiReturn.ErrorNoAccess(c)
			return
		}
		updateField := []string{"IconJson", "Icon", "Title", "Url", "LanUrl", "Description", "OpenMethod", "GroupId"}
		if req.Sort != 0 {
			updateField = append(updateField, "Sort")
		}
		global.Db.Model(&models.ItemIconGroup{}).
			Select(updateField).
			Where("id=?", req.ID).Updates(&req)
	} else {
		// 创建：归属当前用户，自动继承部门
		req.UserId = userInfo.ID
		req.DepartmentId = userInfo.DepartmentId
		global.Db.Create(&req)
	}

	apiReturn.SuccessData(c, req)
}

func (a *ItemIconGroup) GetList(c *gin.Context) {

	userInfo, _ := base.GetCurrentUserInfo(c)
	visitMode := base.GetCurrentVisitMode(c)
	groups := []models.ItemIconGroup{}

	query := global.Db.Order("sort ,created_at")

	if visitMode == base.VISIT_MODE_PUBLIC {
		// 公开访问模式：保持原逻辑，只查公开用户的个人数据
		query = query.Where("user_id=?", userInfo.ID)
	} else {
		// 登录模式：应用部门数据隔离
		query = department.BuildDepartmentScope(query, userInfo.ID, userInfo)
	}

	if err := query.Find(&groups).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 注意：列表接口保持只读，不再自动创建默认分组。
	// 早期实现会在分组为空时自动插入一个名为 "APP" 的默认分组，
	// 导致用户删除该分组后，只要其可见分组集合变为空就会被重新创建，
	// 表现为“删除后分组又出现”。删除行为现在完全由 Deletes 接口控制。
	apiReturn.SuccessListData(c, groups, 0)
}

func (a *ItemIconGroup) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	// 超级管理员可删除任意分组；其他用户只能删除自己的分组
	scopeUserId := userInfo.ID
	if userInfo.Role == department.ROLE_SUPER_ADMIN {
		scopeUserId = 0
	}

	var count int64
	scope := global.Db.Model(&models.ItemIconGroup{})
	if scopeUserId > 0 {
		scope = scope.Where("user_id=?", scopeUserId)
	}
	if err := scope.Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if math.Abs(float64(len(req.Ids))-float64(count)) < 1 {
		apiReturn.ErrorCode(c, 1201, "At least one must be retained", nil)
		return
	}

	txErr := global.Db.Transaction(func(tx *gorm.DB) error {
		mitemIcon := models.ItemIcon{}
		q := tx.Where("id in ?", req.Ids)
		if scopeUserId > 0 {
			q = q.Where("user_id=?", scopeUserId)
		}
		deleteResult := q.Delete(&models.ItemIconGroup{})
		if deleteResult.Error != nil {
			return deleteResult.Error
		}

		// 非管理员只能删除自己拥有的分组。若实际未删除任何行，说明请求的分组不属于当前用户
		// 或不存在（例如共享分组 user_id 为 NULL）。此时不应静默返回成功，否则前端会认为已删除、
		// 刷新后分组又出现（即“删除后又出现”的问题），应明确返回错误。
		if scopeUserId > 0 && deleteResult.RowsAffected == 0 {
			return errGroupNotOwned
		}

		if err := mitemIcon.DeleteByItemIconGroupIds(tx, scopeUserId, req.Ids); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		if errors.Is(txErr, errGroupNotOwned) {
			apiReturn.ErrorCode(c, 1005, "no permission to delete the requested group(s)", nil)
			return
		}
		apiReturn.ErrorDatabase(c, txErr.Error())
		return
	}

	apiReturn.Success(c)
}

// 保存排序
func (a *ItemIconGroup) SaveSort(c *gin.Context) {
	req := commonApiStructs.SortRequest{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	// 超级管理员可排序任意分组；其他用户只能排序自己的分组
	scopeUserId := userInfo.ID
	if userInfo.Role == department.ROLE_SUPER_ADMIN {
		scopeUserId = 0
	}

	transactionErr := global.Db.Transaction(func(tx *gorm.DB) error {
		for _, v := range req.SortItems {
			q := tx.Model(&models.ItemIconGroup{}).Where("id=?", v.Id)
			if scopeUserId > 0 {
				q = q.Where("user_id=?", scopeUserId)
			}
			if err := q.Update("sort", v.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if transactionErr != nil {
		apiReturn.ErrorDatabase(c, transactionErr.Error())
		return
	}

	apiReturn.Success(c)
}
