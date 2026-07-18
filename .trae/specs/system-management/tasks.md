# Sun-Panel 权限表格及系统管理功能模块 - 实现计划

## [ ] Task 1: 数据库模型设计与初始化
- **Priority**: high
- **Depends On**: None
- **Description**: 
  - 设计并创建角色表(Role)、权限表(Permission)、角色权限关联表(RolePermission)、部门表(Department)、备份记录表(Backup)、图标表(Icon)、图标分类表(IconCategory)、图库表(Gallery)等数据库模型
  - 实现权限矩阵初始化脚本，读取 `role_permission_matrix_v2.xlsx` 文件并导入数据库
- **Acceptance Criteria Addressed**: AC-4
- **Test Requirements**:
  - `programmatic` TR-1.1: 所有数据库表成功创建，包含正确的字段和索引
  - `programmatic` TR-1.2: 权限矩阵初始化脚本执行成功，数据正确导入数据库
  - `human-judgement` TR-1.3: 模型设计符合业务需求，关联关系正确

## [ ] Task 2: 后端API - 角色管理接口
- **Priority**: high
- **Depends On**: Task 1
- **Description**: 
  - 实现角色CRUD接口：创建角色、查询角色列表(分页)、更新角色、删除角色
  - 实现角色名称唯一性验证
  - 实现角色删除限制：已分配用户的角色不允许删除
- **Acceptance Criteria Addressed**: AC-6, AC-7
- **Test Requirements**:
  - `programmatic` TR-2.1: POST /api/role/create 创建角色成功，返回角色ID
  - `programmatic` TR-2.2: POST /api/role/list 分页查询角色列表，返回正确的分页数据
  - `programmatic` TR-2.3: POST /api/role/update 更新角色信息成功
  - `programmatic` TR-2.4: POST /api/role/delete 删除已分配用户的角色返回错误
  - `programmatic` TR-2.5: 创建重复名称角色返回错误

## [ ] Task 3: 后端API - 权限配置接口
- **Priority**: high
- **Depends On**: Task 1
- **Description**: 
  - 实现权限矩阵数据获取接口
  - 实现权限配置保存接口，支持批量操作
  - 实现权限缓存刷新机制，确保权限实时生效
- **Acceptance Criteria Addressed**: AC-4, AC-5
- **Test Requirements**:
  - `programmatic` TR-3.1: POST /api/permission/get 获取权限矩阵数据成功
  - `programmatic` TR-3.2: POST /api/permission/save 批量保存权限配置成功
  - `programmatic` TR-3.3: 权限修改后缓存正确刷新

## [ ] Task 4: 后端API - 部门管理接口
- **Priority**: high
- **Depends On**: Task 1
- **Description**: 
  - 实现部门CRUD接口：创建部门、查询部门列表(树形)、更新部门、删除部门
  - 实现部门负责人设置功能
  - 实现部门层级关系管理
- **Acceptance Criteria Addressed**: AC-10, AC-11
- **Test Requirements**:
  - `programmatic` TR-4.1: POST /api/department/create 创建部门成功
  - `programmatic` TR-4.2: POST /api/department/list 返回树形结构部门数据
  - `programmatic` TR-4.3: POST /api/department/update 更新部门信息成功
  - `programmatic` TR-4.4: POST /api/department/delete 删除部门成功

## [ ] Task 5: 后端API - 备份管理接口
- **Priority**: medium
- **Depends On**: Task 1
- **Description**: 
  - 实现备份创建接口，支持数据库备份和全部数据备份两种模式
  - 实现备份恢复接口，包含二次确认逻辑
  - 实现备份删除、导出、导入接口
  - 实现定时备份任务配置接口
- **Acceptance Criteria Addressed**: AC-18, AC-19, AC-20
- **Test Requirements**:
  - `programmatic` TR-5.1: POST /api/backup/create 创建备份成功，返回备份文件路径
  - `programmatic` TR-5.2: POST /api/backup/list 查询备份历史列表成功
  - `programmatic` TR-5.3: POST /api/backup/restore 恢复备份成功
  - `programmatic` TR-5.4: POST /api/backup/schedule 配置定时备份任务成功

## [ ] Task 6: 后端API - 图标管理接口
- **Priority**: medium
- **Depends On**: Task 1
- **Description**: 
  - 实现图标查询接口，支持按ID和名称查询
  - 实现批量获取图标接口，支持按分类批量获取
  - 实现图标分类管理接口
  - 实现图标收藏功能接口
- **Acceptance Criteria Addressed**: AC-21, AC-22
- **Test Requirements**:
  - `programmatic` TR-6.1: POST /api/icon/get 根据ID获取图标成功
  - `programmatic` TR-6.2: POST /api/icon/search 按名称搜索图标成功
  - `programmatic` TR-6.3: POST /api/icon/batchGet 批量获取图标成功
  - `programmatic` TR-6.4: POST /api/icon/favorite 图标收藏功能正常

## [ ] Task 7: 后端API - 图库管理接口
- **Priority**: medium
- **Depends On**: Task 1
- **Description**: 
  - 实现图库图片上传接口
  - 实现图库图片查询接口，支持分类筛选
  - 实现图库图片删除接口，支持批量删除
  - 实现图库图片分类管理接口
- **Acceptance Criteria Addressed**: AC-23, AC-24
- **Test Requirements**:
  - `programmatic` TR-7.1: POST /api/gallery/upload 上传图片成功
  - `programmatic` TR-7.2: POST /api/gallery/list 查询图库列表成功
  - `programmatic` TR-7.3: POST /api/gallery/delete 批量删除图片成功

## [ ] Task 8: 后端API - 系统设置接口(Logo/背景)
- **Priority**: high
- **Depends On**: Task 1
- **Description**: 
  - 实现Logo设置接口，支持上传、尺寸配置、CDN路径配置
  - 实现登录页背景设置接口，支持上传、URL配置、显示方式设置
  - 实现预设背景图片管理接口
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3, AC-12, AC-13, AC-14
- **Test Requirements**:
  - `programmatic` TR-8.1: POST /api/system/logo/set 设置Logo成功
  - `programmatic` TR-8.2: POST /api/system/logo/get 获取Logo配置成功
  - `programmatic` TR-8.3: POST /api/system/background/set 设置背景成功
  - `programmatic` TR-8.4: POST /api/system/background/get 获取背景配置成功

## [ ] Task 9: 前端 - 角色权限管理页面
- **Priority**: high
- **Depends On**: Task 2, Task 3
- **Description**: 
  - 开发角色选择组件，支持搜索过滤
  - 开发权限矩阵展示组件，支持模块开关和权限点勾选
  - 开发角色新增/编辑/删除弹窗组件
  - 实现全选/反选功能
- **Acceptance Criteria Addressed**: AC-5, AC-6, AC-7, AC-8, AC-9
- **Test Requirements**:
  - `human-judgement` TR-9.1: 角色选择下拉框支持搜索过滤功能
  - `human-judgement` TR-9.2: 权限矩阵展示清晰，模块开关控制正常
  - `human-judgement` TR-9.3: 角色新增表单验证正常
  - `human-judgement` TR-9.4: 删除角色二次确认弹窗正常

## [ ] Task 10: 前端 - 部门管理页面
- **Priority**: high
- **Depends On**: Task 4
- **Description**: 
  - 开发部门树形结构展示组件
  - 开发部门新增/编辑弹窗组件
  - 开发部门负责人选择组件，支持模糊搜索
- **Acceptance Criteria Addressed**: AC-10, AC-11
- **Test Requirements**:
  - `human-judgement` TR-10.1: 部门树形结构展示清晰，层级关系正确
  - `human-judgement` TR-10.2: 部门新增/编辑表单验证正常
  - `human-judgement` TR-10.3: 负责人选择组件支持模糊搜索

## [ ] Task 11: 前端 - 系统Logo设置页面
- **Priority**: high
- **Depends On**: Task 8
- **Description**: 
  - 开发Logo上传组件，支持PNG/JPG/JPEG格式，限制5MB
  - 开发Logo预览组件，保持与实际展示区域比例一致
  - 开发Logo尺寸调整控件(50px-200px)
  - 修改登录页面，在文字上方居中显示Logo
- **Acceptance Criteria Addressed**: AC-1, AC-2, AC-3
- **Test Requirements**:
  - `human-judgement` TR-11.1: Logo上传组件格式和大小限制正常
  - `human-judgement` TR-11.2: Logo预览效果与实际展示一致
  - `human-judgement` TR-11.3: 尺寸调整控件范围正确，实时预览正常
  - `human-judgement` TR-11.4: 登录页面Logo居中显示，适配不同屏幕

## [ ] Task 12: 前端 - 登录页背景设置页面
- **Priority**: medium
- **Depends On**: Task 8
- **Description**: 
  - 开发背景上传组件，支持上传和URL配置
  - 开发预设背景选择组件，提供至少5种预设背景
  - 开发背景显示方式选择组件(平铺/拉伸/居中)
  - 实现背景效果实时预览
- **Acceptance Criteria Addressed**: AC-12, AC-13, AC-14
- **Test Requirements**:
  - `human-judgement` TR-12.1: 背景上传和URL配置功能正常
  - `human-judgement` TR-12.2: 预设背景库至少包含5种背景
  - `human-judgement` TR-12.3: 背景显示方式设置正常，预览效果正确

## [ ] Task 13: 前端 - 文件上传管理组件
- **Priority**: medium
- **Depends On**: Task 7
- **Description**: 
  - 开发多文件上传组件，支持拖拽上传
  - 开发上传进度条组件，显示上传百分比
  - 实现上传状态展示(上传中/成功/失败)
  - 实现失败重试和取消上传功能
- **Acceptance Criteria Addressed**: AC-15, AC-16, AC-17
- **Test Requirements**:
  - `human-judgement` TR-13.1: 多文件拖拽上传功能正常
  - `human-judgement` TR-13.2: 上传进度条准确显示百分比
  - `human-judgement` TR-13.3: 上传失败后重试功能正常
  - `human-judgement` TR-13.4: 取消上传功能正常

## [ ] Task 14: 前端 - 备份管理页面
- **Priority**: medium
- **Depends On**: Task 5
- **Description**: 
  - 开发备份模式选择组件(数据库备份/全部数据备份)
  - 开发备份历史记录列表组件
  - 开发备份创建/恢复/删除/导出/导入功能
  - 开发定时备份任务配置组件
- **Acceptance Criteria Addressed**: AC-18, AC-19, AC-20
- **Test Requirements**:
  - `human-judgement` TR-14.1: 备份模式选择清晰，区别明确
  - `human-judgement` TR-14.2: 备份历史记录查询正常
  - `human-judgement` TR-14.3: 恢复操作二次确认正常
  - `human-judgement` TR-14.4: 定时备份任务配置正常

## [ ] Task 15: 前端 - 图标管理页面
- **Priority**: medium
- **Depends On**: Task 6
- **Description**: 
  - 开发图标列表展示组件，支持按分类筛选
  - 开发图标搜索组件，支持按名称搜索
  - 开发图标收藏功能组件
  - 开发图标批量操作组件
- **Acceptance Criteria Addressed**: AC-21, AC-22
- **Test Requirements**:
  - `human-judgement` TR-15.1: 图标列表按分类展示清晰
  - `human-judgement` TR-15.2: 图标搜索功能正常
  - `human-judgement` TR-15.3: 图标收藏功能正常
  - `human-judgement` TR-15.4: 批量操作功能正常

## [ ] Task 16: 前端 - 图库管理页面
- **Priority**: medium
- **Depends On**: Task 7, Task 13
- **Description**: 
  - 在系统管理菜单下新增图库菜单，包含壁纸管理和图标管理子菜单
  - 开发图库图片上传组件(点击上传/拖拽上传)
  - 开发图库图片预览和管理组件
  - 开发图片分类管理组件
- **Acceptance Criteria Addressed**: AC-23, AC-24
- **Test Requirements**:
  - `human-judgement` TR-16.1: 系统管理菜单下正确显示图库菜单及子菜单
  - `human-judgement` TR-16.2: 图片上传支持点击和拖拽方式
  - `human-judgement` TR-16.3: 图片预览和删除功能正常
  - `human-judgement` TR-16.4: 图片分类管理功能正常

## [ ] Task 17: 前端 - API接口封装
- **Priority**: high
- **Depends On**: Task 2, Task 3, Task 4, Task 5, Task 6, Task 7, Task 8
- **Description**: 
  - 封装角色管理API接口
  - 封装权限配置API接口
  - 封装部门管理API接口
  - 封装备份管理API接口
  - 封装图标管理API接口
  - 封装图库管理API接口
  - 封装系统设置API接口(Logo/背景)
- **Acceptance Criteria Addressed**: 所有AC
- **Test Requirements**:
  - `programmatic` TR-17.1: 所有API接口封装正确，调用正常
  - `human-judgement` TR-17.2: API接口命名规范，结构清晰

## [ ] Task 18: 路由配置与权限拦截
- **Priority**: high
- **Depends On**: Task 9, Task 10, Task 11, Task 12, Task 14, Task 15, Task 16
- **Description**: 
  - 配置系统管理模块路由
  - 实现管理员权限拦截，仅管理员可访问系统管理模块
  - 配置各功能页面路由
- **Acceptance Criteria Addressed**: 所有AC
- **Test Requirements**:
  - `programmatic` TR-18.1: 系统管理路由配置正确，页面可正常访问
  - `programmatic` TR-18.2: 非管理员访问系统管理模块被正确拦截
  - `human-judgement` TR-18.3: 路由导航正常，页面跳转流畅

## [ ] Task 19: 国际化与样式优化
- **Priority**: low
- **Depends On**: 所有前端任务
- **Description**: 
  - 添加系统管理模块相关的中英文翻译
  - 优化各页面样式，保持与现有系统风格一致
  - 确保响应式设计适配不同屏幕尺寸
- **Acceptance Criteria Addressed**: 所有AC
- **Test Requirements**:
  - `human-judgement` TR-19.1: 中英文切换正常，翻译准确
  - `human-judgement` TR-19.2: 页面样式与现有系统风格一致
  - `human-judgement` TR-19.3: 响应式设计适配不同屏幕尺寸

## [ ] Task 20: 测试与验证
- **Priority**: high
- **Depends On**: 所有任务
- **Description**: 
  - 执行后端API单元测试
  - 执行前端功能测试
  - 验证所有验收标准
  - 修复发现的问题
- **Acceptance Criteria Addressed**: 所有AC
- **Test Requirements**:
  - `programmatic` TR-20.1: 后端API测试用例全部通过
  - `human-judgement` TR-20.2: 前端功能测试验证所有验收标准
  - `human-judgement` TR-20.3: 系统运行稳定，无明显bug