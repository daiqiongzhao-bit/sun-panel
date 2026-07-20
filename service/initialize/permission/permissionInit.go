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

	// 为部门管理员分配默认权限
	assignDeptAdminPermissions()

	return nil
}

// assignDeptAdminPermissions 为部门管理员角色分配默认权限
func assignDeptAdminPermissions() error {
	deptAdminRole := models.Role{}
	if err := models.Db.Where("name=?", "部门管理员").First(&deptAdminRole).Error; err != nil {
		return nil // 角色不存在则跳过
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

	// 清理旧权限并分配新权限
	tx := models.Db.Begin()
	if err := tx.Delete(&models.RolePermission{}, "role_id=?", deptAdminRole.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
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

	// 为"管理员"角色补充管理相关权限（保留已有手动调整，仅补齐缺失项）
	mgrRole := models.Role{}
	if err := tx.Where("name=?", "管理员").First(&mgrRole).Error; err == nil {
		var mgrPermissions []models.Permission
		if err := tx.Where("module_code != ? OR module_code IS NULL", "permission").Find(&mgrPermissions).Error; err == nil {
			for _, p := range mgrPermissions {
				var cnt int64
				tx.Model(&models.RolePermission{}).Where("role_id=? AND permission_id=?", mgrRole.ID, p.ID).Count(&cnt)
				if cnt == 0 {
					tx.Create(&models.RolePermission{RoleId: mgrRole.ID, PermissionId: p.ID})
				}
			}
		}
	}

	// 为"普通用户"角色分配个人功能默认权限（补充缺失项，保留已有手动调整）
	normalRole := models.Role{}
	if err := tx.Where("name=?", "普通用户").First(&normalRole).Error; err == nil {
		normalPermCodes := []string{
			"dashboard:view",
			"style:view", "style:edit",
			"group:view", "group:create", "group:edit", "group:delete",
			"file:view", "file:upload", "file:delete",
			"icon:view", "icon:upload", "icon:delete",
			"sticky_note:view", "sticky_note:create", "sticky_note:edit", "sticky_note:delete",
			"import_export:export", "import_export:import",
			"paste:view", "paste:create",
			"monitor:view",
		}
		var normalPerms []models.Permission
		if err := tx.Where("permission_id IN ?", normalPermCodes).Find(&normalPerms).Error; err == nil {
			for _, p := range normalPerms {
				var cnt int64
				tx.Model(&models.RolePermission{}).Where("role_id=? AND permission_id=?", normalRole.ID, p.ID).Count(&cnt)
				if cnt == 0 {
					tx.Create(&models.RolePermission{RoleId: normalRole.ID, PermissionId: p.ID})
				}
			}
		}
	}

	tx.Commit()

	// 为部门管理员分配默认权限
	assignDeptAdminPermissions()

	return nil
}
