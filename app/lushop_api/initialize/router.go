package initialize

import (
	"lushopapi/middlewares"
	"lushopapi/router"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 注册pprof路由
	// 将标准库pprof的Handler注册到Gin路由，路径前缀为/debug/pprof
	Router.Any("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	Router.LoadHTMLFiles("template/weboss.html")
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	Router.StaticFS("/static", http.Dir("static"))
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	Router.GET("/", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "weboss.html", gin.H{
			"title": "posts/index",
		})
	})
	//配置跨域
	Router.Use(middlewares.Cors())

	router.InitUserRouter(Router)   //用户相关
	router.InitUserOpRouter(Router) //用户操作相关
	// router.InitBaseRouter(Router)   //验证码相关
	router.InitGoodsRouter(Router)  //商品相关
	router.InitOrderRouter(Router)  //订单相关
	router.InitOssRouter(Router)    //oss相关
	router.InitOAuth2Router(Router) // oauth2相关
	return Router
}
