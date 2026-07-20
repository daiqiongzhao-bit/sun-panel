package panel

import "github.com/gin-gonic/gin"

func Init(routerGroup *gin.RouterGroup) {
	InitItemIconGroup(routerGroup)
	InitItemIcon(routerGroup)
	InitUserConfig(routerGroup)
	InitUsersRouter(routerGroup)
	InitStickyNoteRouter(routerGroup)
	InitPasteBinRouter(routerGroup)
	InitSearchRouter(routerGroup)
}
