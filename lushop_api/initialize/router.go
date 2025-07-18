package initialize

import (
	"lushopapi/middlewares"
	"lushopapi/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		return
	})
	//配置跨域
	Router.Use(middlewares.Cors())
	//ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(Router)   //用户相关
	router.InitUserOpRouter(Router) //用户操作相关
	router.InitBaseRouter(Router)   //验证码相关
	router.InitGoodsRouter(Router)  //商品相关
	router.InitOrderRouter(Router)  //订单相关
	router.InitOssRouter(Router)    //oss相关
	return Router
}
