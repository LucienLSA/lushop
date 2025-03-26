package initialize

import (
	"lushopapi/user_web/global"
	"lushopapi/user_web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// fmt.Println(global.ServerConfig.UserSrvInfo.Version)
	ApiGroup := Router.Group("/u/" + global.ServerConfig.UserSrvInfo.Version)
	router.InitUserRouter(ApiGroup)
	return Router
}
