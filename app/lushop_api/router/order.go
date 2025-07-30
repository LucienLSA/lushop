package router

import (
	"lushopapi/api/order"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("/o")
	OrderRouter := ApiGroup.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.OrderList)       // 订单列表
		OrderRouter.POST("", order.OrderCreate)    //新建订单
		OrderRouter.GET("/:id", order.OrderDetail) //订单详情
	}
	PayRouter := ApiGroup.Group("pay").Use(middlewares.JWTAuth())
	{
		PayRouter.POST("/notify", order.Notify) //支付回调通知
	}
	GoodsRouter := ApiGroup.Group("shopcart").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", order.ShopList)          //购物车列表
		GoodsRouter.DELETE("/:id", order.ShopDelete) //删除条目
		GoodsRouter.POST("", order.ShopCreate)       //添加商品购物车
		GoodsRouter.PUT("/:id", order.ShopUpdate)    //修改商品购物车
	}
}
