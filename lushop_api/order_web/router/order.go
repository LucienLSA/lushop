package router

import (
	"orderweb/api/order"
	"orderweb/api/pay"
	"orderweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("/list", order.List)         // 获取订单列表
		OrderRouter.POST("/new", order.New)          // 新建订单
		OrderRouter.GET("/detail/:id", order.Detail) // 订单详情
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("/alipay/notify", pay.Notify)
	}
}
