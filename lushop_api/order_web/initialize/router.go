package initialize

import (
	"net/http"
	"orderweb/global"
	"orderweb/middlewares"
	"orderweb/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	Router.Use(middlewares.Cors())
	// fmt.Println(global.ServerConfig.UserSrvInfo.Version)
	ApiGroup := Router.Group("/o/" + global.ServerConfig.OrderSrvInfo.Version)
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)
	return Router
}
