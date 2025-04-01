package initialize

import (
	"lushopapi/user_web/global"
	"lushopapi/user_web/middlewares"
	"lushopapi/user_web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	// fmt.Println(global.ServerConfig.UserSrvInfo.Version)
	ApiGroup := Router.Group("/u/" + global.ServerConfig.UserSrvInfo.Version)
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	return Router
}
