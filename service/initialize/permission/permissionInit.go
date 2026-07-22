package permission

import (
	"fmt"
	"sun-panel/models"

	"github.com/xuri/excelize/v2"
)

func InitPermissionMatrix(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取Excel数据失败: %v", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("Excel文件数据行数不足")
	}

	headers := rows[0]
	roleNames := headers[2:]

	modules := make(map[string]string)
	permissions := make([]models.Permission, 0)
	rolePermissionMap := make(map[string][]string)

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 3 {
			continue
		}

		moduleCode := row[0]
		moduleName := row[1]
		permissionId := row[0] + ":" + row[1]
		permissionName := row[1]

		if moduleCode != "" && moduleName != "" {
			modules[moduleCode] = moduleName
		}

		permission := models.Permission{
			ModuleCode:     moduleCode,
			ModuleName:     moduleName,
			PermissionId:   permissionId,
			PermissionName: permissionName,
			Status:         1,
		}
		permissions = append(permissions, permission)

		for j := 2; j < len(row) && j-2 < len(roleNames); j++ {
			roleName := roleNames[j-2]
			value := row[j]
			if value == "1" || value == "是" || value == "Y" || value == "y" {
				if _, ok := rolePermissionMap[roleName]; !ok {
					rolePermissionMap[roleName] = make([]string, 0)
				}
				rolePermissionMap[roleName] = append(rolePermissionMap[roleName], permissionId)
			}
		}
	}

	tx := models.Db.Begin()

	for _, permission := range permissions {
		var exists models.Permission
		if err := tx.Where("permission_id=?", permission.PermissionId).First(&exists).Error; err != nil {
			if err == models.Db.Error {
				if err := tx.Create(&permission).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("创建权限失败: %v", err)
				}
			}
		}
	}

	for roleName := range rolePermissionMap {
		var role models.Role
		if err := tx.Where("name=?", roleName).First(&role).Error; err != nil {
			role = models.Role{
				Name:        roleName,
				Description: roleName + "角色",
				Status:      1,
			}
			if err := tx.Create(&role).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("创建角色失败: %v", err)
			}
		}
	}

	tx.Commit()

	// 为部门管理员分配默认权限（仅当没有任何权限时）
	assignDeptAdminPermissionsIfEmpty()

	return nil
}

// assignDeptAdminPermissionsIfEmpty 仅在部门管理员角色没有任何权限时分配默认权限
// 避免每次启动覆盖管理员在角色权限页做的自定义调整
func assignDeptAdminPermissionsIfEmpty() error {
	deptAdminRole := models.Role{}
	if err := models.Db.Where("name=?", "部门管理员").First(&deptAdminRole).Error; err != nil {
		return nil // 角色不存在则跳过
	}

	// 如果已有权限，不再覆盖
	var cnt int64
	if err := models.Db.Model(&models.RolePermission{}).Where("role_id=?", deptAdminRole.ID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	// 部门管理员可拥有的权限ID列表
	allowedPermissionIds := []string{
		"dashboard:view",
		"user:view",
		"user:edit",
		"user:create",
		"department:view",
		"department:edit",
		"department:create",
		"system:view",
		"setting:view",
		"monitor:view",
	}

	var permissions []models.Permission
	if err := models.Db.Where("permission_id IN ?", allowedPermissionIds).Find(&permissions).Error; err != nil {
		return err
	}

	// 当前没有任何权限，直接分配默认权限
	tx := models.Db.Begin()
	for _, p := range permissions {
		rp := models.RolePermission{
			RoleId:       deptAdminRole.ID,
			PermissionId: p.ID,
		}
		if err := tx.Create(&rp).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func InitDefaultPermissions() error {
	defaultModules := []struct {
		Code string
		Name string
	}{
		{"system", "系统管理"},
		{"dashboard", "仪表盘"},
		{"user", "用户管理"},
		{"role", "角色管理"},
		{"permission", "权限管理"},
		{"department", "部门管理"},
		{"backup", "备份管理"},
		{"gallery", "图库管理"},
		{"icon", "图标管理"},
		{"setting", "系统设置"},
		{"sticky_note", "便签管理"},
		{"operation_log", "操作日志"},
		{"login_log", "登录日志"},
		{"notice", "通知公告"},
		{"file", "文件管理"},
		{"monitor", "系统监控"},
		{"style", "样式设置"},
		{"import_export", "导入导出"},
		{"task", "任务管理"},
		{"paste", "粘贴板"},
		{"about", "关于"},
		{"group", "分组管理"},
	}

	defaultPermissions := []struct {
		ModuleCode     string
		PermissionId   string
		PermissionName string
	}{
		{"system", "system:view", "查看"},
		{"system", "system:edit", "编辑"},
		{"system", "system:create", "创建"},
		{"system", "system:delete", "删除"},
		{"dashboard", "dashboard:view", "查看"},
		{"user", "user:view", "查看"},
		{"user", "user:edit", "编辑"},
		{"user", "user:create", "创建"},
		{"user", "user:delete", "删除"},
		{"role", "role:view", "查看"},
		{"role", "role:edit", "编辑"},
		{"role", "role:create", "创建"},
		{"role", "role:delete", "删除"},
		{"permission", "permission:view", "查看"},
		{"permission", "permission:edit", "编辑"},
		{"department", "department:view", "查看"},
		{"department", "department:edit", "编辑"},
		{"department", "department:create", "创建"},
		{"department", "department:delete", "删除"},
		{"backup", "backup:view", "查看"},
		{"backup", "backup:create", "创建"},
		{"backup", "backup:restore", "恢复"},
		{"backup", "backup:delete", "删除"},
		{"gallery", "gallery:view", "查看"},
		{"gallery", "gallery:upload", "上传"},
		{"gallery", "gallery:delete", "删除"},
		{"icon", "icon:view", "查看"},
		{"icon", "icon:create", "创建"},
		{"icon", "icon:edit", "编辑"},
		{"icon", "icon:upload", "上传"},
		{"icon", "icon:delete", "删除"},
		{"setting", "setting:view", "查看"},
		{"setting", "setting:edit", "编辑"},
		{"sticky_note", "sticky_note:view", "查看"},
		{"sticky_note", "sticky_note:create", "创建"},
		{"sticky_note", "sticky_note:edit", "编辑"},
		{"sticky_note", "sticky_note:delete", "删除"},
		{"operation_log", "operation_log:view", "查看"},
		{"operation_log", "operation_log:clear", "清空"},
		{"login_log", "login_log:view", "查看"},
		{"login_log", "login_log:clear", "清空"},
		{"notice", "notice:view", "查看"},
		{"notice", "notice:create", "创建"},
		{"notice", "notice:edit", "编辑"},
		{"notice", "notice:delete", "删除"},
		{"file", "file:view", "查看"},
		{"file", "file:upload", "上传"},
		{"file", "file:delete", "删除"},
		{"monitor", "monitor:view", "查看"},
		{"style", "style:view", "查看"},
		{"style", "style:edit", "编辑"},
		{"import_export", "import_export:export", "导出"},
		{"import_export", "import_export:import", "导入"},
		{"task", "task:view", "查看"},
		{"task", "task:create", "创建"},
		{"task", "task:edit", "编辑"},
		{"task", "task:delete", "删除"},
		{"paste", "paste:view", "查看"},
		{"paste", "paste:create", "创建"},
		{"paste", "paste:edit", "编辑"},
		{"paste", "paste:delete", "删除"},
		{"about", "about:view", "查看"},
		{"group", "group:view", "查看"},
		{"group", "group:create", "创建"},
		{"group", "group:edit", "编辑"},
		{"group", "group:delete", "删除"},
	}

	defaultRoles := []struct {
		Name        string
		Description string
	}{
		{"超级管理员", "拥有系统所有权限"},
		{"管理员", "拥有系统管理权限"},
		{"部门管理员", "管理本部门及子部门的数据和用户"},
		{"普通用户", "拥有基础访问权限"},
	}

	tx := models.Db.Begin()

	for _, module := range defaultModules {
		var count int64
		tx.Model(&models.Permission{}).Where("module_code=?", module.Code).Count(&count)
		if count == 0 {
			for _, perm := range defaultPermissions {
				if perm.ModuleCode == module.Code {
					var exists models.Permission
					if err := tx.Where("permission_id=?", perm.PermissionId).First(&exists).Error; err != nil {
						permission := models.Permission{
							ModuleCode:     module.Code,
							ModuleName:     module.Name,
							PermissionId:   perm.PermissionId,
							PermissionName: perm.PermissionName,
							Status:         1,
						}
						if err := tx.Create(&permission).Error; err != nil {
							tx.Rollback()
							return fmt.Errorf("创建权限失败: %v", err)
						}
					}
				}
			}
		}
	}

	// 首次创建系统内置角色时赋予默认权限；已存在的角色不再覆盖，确保管理员自定义持久化
	defaultRolePermissions := map[string][]string{
		"普通用户": {
			"dashboard:view",
			"style:view", "style:edit",
			"group:view",
			"file:view",
			"icon:view",
			"sticky_note:view", "sticky_note:create", "sticky_note:edit", "sticky_note:delete",
			"paste:view", "paste:create",
			"monitor:view",
			"about:view",
		},
		"管理员": nil, // nil 表示全量权限
	}

	createdRoleIds := make(map[string]uint)
	for _, role := range defaultRoles {
		var exists models.Role
		if err := tx.Where("name=?", role.Name).First(&exists).Error; err != nil {
			newRole := models.Role{
				Name:        role.Name,
				Description: role.Description,
				Status:      1,
			}
			if err := tx.Create(&newRole).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("创建角色失败: %v", err)
			}
			createdRoleIds[role.Name] = newRole.ID
		}
	}

	// 为本次新创建的角色分配默认权限
	var allPermissions []models.Permission
	if err := tx.Find(&allPermissions).Error; err == nil {
		permIdByCode := make(map[string]uint)
		for _, p := range allPermissions {
			permIdByCode[p.PermissionId] = p.ID
		}
		for roleName, roleId := range createdRoleIds {
			permCodes := defaultRolePermissions[roleName]
			if permCodes == nil {
				// 管理员：分配全部权限（除权限管理模块，保留给超级管理员）
				for _, p := range allPermissions {
					if p.ModuleCode != "permission" {
						tx.Create(&models.RolePermission{RoleId: roleId, PermissionId: p.ID})
					}
				}
			} else {
				for _, code := range permCodes {
					if pid, ok := permIdByCode[code]; ok {
						tx.Create(&models.RolePermission{RoleId: roleId, PermissionId: pid})
					}
				}
			}
		}
	}

	adminRole := models.Role{}
	if err := tx.Where("name=?", "超级管理员").First(&adminRole).Error; err == nil {
		var allPermissions []models.Permission
		if err := tx.Find(&allPermissions).Error; err == nil {
			var permIds []uint
			for _, p := range allPermissions {
				permIds = append(permIds, p.ID)
			}
			if err := tx.Delete(&models.RolePermission{}, "role_id=?", adminRole.ID).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("清理角色权限失败: %v", err)
			}
			for _, pid := range permIds {
				rp := models.RolePermission{
					RoleId:       adminRole.ID,
					PermissionId: pid,
				}
				if err := tx.Create(&rp).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("分配权限失败: %v", err)
				}
			}
		}
	}

	// 注意：不再自动给"管理员"和"普通用户"补充默认权限，避免覆盖管理员在角色权限页做的自定义调整。
	// 仅当角色完全不存在的系统首次初始化时，由管理员手动在"角色权限"页面配置。

	tx.Commit()

	// 为部门管理员分配默认权限（仅当该角色没有任何权限时，避免覆盖手动配置）
	assignDeptAdminPermissionsIfEmpty()

	return nil
}
