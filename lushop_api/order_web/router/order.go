package router

import (
	"lushopapi/order_web/api/order"
	"lushopapi/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order")
	{
		OrderRouter.GET("/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), order.List) // 获取订单列表
		OrderRouter.POST("/new", middlewares.JWTAuth(), order.New)                             // 新建订单
		OrderRouter.GET("/detail/:id", order.Detail)                                           // 订单详情
	}
}
