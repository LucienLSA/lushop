package initialize

import (
	"lushopapi/goods_web/global"
	"lushopapi/goods_web/middlewares"
	"lushopapi/goods_web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	// fmt.Println(global.ServerConfig.UserSrvInfo.Version)
	ApiGroup := Router.Group("/g/" + global.ServerConfig.GoodsSrvInfo.Version)
	router.InitGoodsRouter(ApiGroup)
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	return Router
}
