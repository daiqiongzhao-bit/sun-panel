package system

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitAbout(routerGroup)
	InitLogin(routerGroup)
	InitUserRouter(routerGroup)
	InitFileRouter(routerGroup)
	InitNoticeRouter(routerGroup)
	InitNoticeManageRouter(routerGroup)
	InitLogRouter(routerGroup)
	InitJobRouter(routerGroup)
	InitModuleConfigRouter(routerGroup)
	InitMonitorRouter(routerGroup)
	InitRoleRouter(routerGroup)
	InitPermissionRouter(routerGroup)
	InitDepartmentRouter(routerGroup)
	InitSystemSettingRouter(routerGroup)
	InitBackupRouter(routerGroup)
	InitIconRouter(routerGroup)
	InitGalleryRouter(routerGroup)
	InitImportExportRouter(routerGroup)
}
