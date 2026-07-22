package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/department"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 初始化一个独立的内存 SQLite 数据库并迁移所需模型，
// 将其挂载到 global.Db 供本包的接口进行回归测试。
func setupTestDB(t *testing.T) {
	t.Helper()
	// 每次测试使用唯一的内存库名，避免测试之间相互影响。
	dsn := fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.ItemIconGroup{}, &models.ItemIcon{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	global.Db = db
}

// newContext 构造一个带 POST body、userInfo 与 visitMode 的 gin 测试上下文。
func newContext(t *testing.T, body interface{}, userInfo models.User, visitMode int) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader([]byte("{}"))
	}
	c.Request, _ = http.NewRequest(http.MethodPost, "/", reader)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userInfo", userInfo)
	c.Set(base.GIN_GET_VISIT_MODE, visitMode)
	return c, w
}

// responseCode 解析接口返回的 code 字段。
func responseCode(t *testing.T, w *httptest.ResponseRecorder) int {
	t.Helper()
	var resp struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp: %v (body=%s)", err, w.Body.String())
	}
	return resp.Code
}

// groupTitlesInDB 直接读取数据库中某用户可见分组（超级管理员看全部）的标题。
func groupTitlesInDB(t *testing.T, user models.User) []string {
	t.Helper()
	var groups []models.ItemIconGroup
	query := global.Db.Order("sort, created_at")
	if user.Role != department.ROLE_SUPER_ADMIN {
		query = query.Where("user_id = ? OR department_id IN ?", user.ID, []uint{})
	}
	if err := query.Find(&groups).Error; err != nil {
		t.Fatalf("query groups: %v", err)
	}
	titles := make([]string, 0, len(groups))
	for _, g := range groups {
		titles = append(titles, g.Title)
	}
	return titles
}

// ---------------------------------------------------------------------------
// 回归测试 1（核心）：空分组集合下，GetList 不应再自动创建 "APP" 默认分组。
// 该场景正是“删除分组后又出现”的根因：旧实现会在分组为空时插入名为 APP 的分组。
// ---------------------------------------------------------------------------
func TestGetListDoesNotAutoCreateAppGroup(t *testing.T) {
	setupTestDB(t)
	admin := models.User{Role: department.ROLE_SUPER_ADMIN}
	admin.ID = 1
	if err := global.Db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}

	// 超级管理员当前没有任何分组
	c, w := newContext(t, nil, admin, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).GetList(c)
	if responseCode(t, w) != 0 {
		t.Fatalf("GetList should succeed, code=%d", responseCode(t, w))
	}

	titles := groupTitlesInDB(t, admin)
	for _, title := range titles {
		if title == "APP" {
			t.Fatalf("GetList 自动创建了默认 APP 分组（回归失败）：%v", titles)
		}
	}
	if len(titles) != 0 {
		t.Fatalf("空分组集合下 GetList 不应创建任何分组，实际：%v", titles)
	}
}

// ---------------------------------------------------------------------------
// 回归测试 2：超级管理员删除 APP 分组后，刷新不再出现，其余分组保持完整。
// ---------------------------------------------------------------------------
func TestDeleteAppGroupStaysDeleted(t *testing.T) {
	setupTestDB(t)
	admin := models.User{Role: department.ROLE_SUPER_ADMIN}
	admin.ID = 1
	if err := global.Db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	groups := []models.ItemIconGroup{
		{Title: "APP", UserId: 1},
		{Title: "飞牛", UserId: 1},
		{Title: "系统状态", UserId: 1},
	}
	if err := global.Db.Create(&groups).Error; err != nil {
		t.Fatalf("create groups: %v", err)
	}

	// 前置：GetList 能看到 APP
	c1, _ := newContext(t, nil, admin, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).GetList(c1)
	titles1 := groupTitlesInDB(t, admin)
	if !contains(titles1, "APP") {
		t.Fatalf("前置条件失败：未找到 APP 分组，%v", titles1)
	}

	// 取出 APP 分组 ID 并删除
	var appGroup models.ItemIconGroup
	if err := global.Db.First(&appGroup, "title = ? AND user_id = ?", "APP", 1).Error; err != nil {
		t.Fatalf("find APP group: %v", err)
	}
	c2, w2 := newContext(t, commonApiStructs.RequestDeleteIds[uint]{Ids: []uint{appGroup.ID}}, admin, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).Deletes(c2)
	if code := responseCode(t, w2); code != 0 {
		t.Fatalf("Deletes 应成功，code=%d", code)
	}

	// 后置：GetList 不再包含 APP，但飞牛/系统状态 仍在
	c3, _ := newContext(t, nil, admin, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).GetList(c3)
	titles3 := groupTitlesInDB(t, admin)
	if contains(titles3, "APP") {
		t.Fatalf("删除后 APP 分组仍然存在（回归失败）：%v", titles3)
	}
	if !contains(titles3, "飞牛") || !contains(titles3, "系统状态") {
		t.Fatalf("删除 APP 后其余分组不应丢失：%v", titles3)
	}
}

// ---------------------------------------------------------------------------
// 回归测试 3：非管理员尝试删除共享分组（user_id 为 NULL）时，
// 不应静默返回成功（旧实现会误删 0 行却返回成功，导致“删除后又出现”）。
// 应返回错误，且分组保持不变；删除自己拥有的分组则应成功。
// ---------------------------------------------------------------------------
func TestNonAdminCannotDeleteSharedGroupSilently(t *testing.T) {
	setupTestDB(t)
	regular := models.User{Role: 2}
	regular.ID = 2
	if err := global.Db.Create(&regular).Error; err != nil {
		t.Fatalf("create regular user: %v", err)
	}
	sharedApp := models.ItemIconGroup{Title: "APP", UserId: 0, DepartmentId: 0} // 共享分组
	ownedGroup := models.ItemIconGroup{Title: "我的分组", UserId: 2}
	ownedGroup2 := models.ItemIconGroup{Title: "我的分组2", UserId: 2}
	if err := global.Db.Create(&[]models.ItemIconGroup{sharedApp, ownedGroup, ownedGroup2}).Error; err != nil {
		t.Fatalf("create groups: %v", err)
	}

	// 取回两个分组的 ID
	var shared, owned models.ItemIconGroup
	if err := global.Db.First(&shared, "title = ? AND user_id = ?", "APP", 0).Error; err != nil {
		t.Fatalf("find shared: %v", err)
	}
	if err := global.Db.First(&owned, "title = ? AND user_id = ?", "我的分组", 2).Error; err != nil {
		t.Fatalf("find owned: %v", err)
	}

	// 删除共享分组（非管理员）→ 应返回错误，且分组不被删除
	c1, w1 := newContext(t, commonApiStructs.RequestDeleteIds[uint]{Ids: []uint{shared.ID}}, regular, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).Deletes(c1)
	if code := responseCode(t, w1); code == 0 {
		t.Fatalf("非管理员删除共享分组应失败，但返回了成功 code=0")
	}
	if !groupExists(t, shared.ID) {
		t.Fatalf("非管理员删除共享分组不应真正删除该分组")
	}

	// 删除自己拥有的分组（非管理员）→ 应成功，且分组被删除
	c2, w2 := newContext(t, commonApiStructs.RequestDeleteIds[uint]{Ids: []uint{owned.ID}}, regular, base.VISIT_MODE_LOGIN)
	(&ItemIconGroup{}).Deletes(c2)
	if code := responseCode(t, w2); code != 0 {
		t.Fatalf("非管理员删除自有分组应成功，code=%d", code)
	}
	if groupExists(t, owned.ID) {
		t.Fatalf("非管理员删除自有分组后该分组应已不存在")
	}
}

// groupExists 检查指定 ID 的分组是否仍存在。
func groupExists(t *testing.T, id uint) bool {
	t.Helper()
	var count int64
	if err := global.Db.Model(&models.ItemIconGroup{}).Where("id = ?", id).Count(&count).Error; err != nil {
		t.Fatalf("count group: %v", err)
	}
	return count > 0
}

// contains 判断字符串切片是否包含目标值。
func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
