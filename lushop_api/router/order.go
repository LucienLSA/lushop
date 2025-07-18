package router

import (
	"lushopapi/api/order"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.Engine) {
	//BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.OrderList)       // 订单列表
		OrderRouter.POST("", order.OrderCreate)    //新建订单
		OrderRouter.GET("/:id", order.OrderDetail) //订单详情
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("/alipay/notify", order.Notify) //支付回调
	}
	GoodsRouter := Router.Group("shopcart").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", order.ShopList)          //购物车列表
		GoodsRouter.DELETE("/:id", order.ShopDelete) //删除条目
		GoodsRouter.POST("/", order.ShopCreate)      //添加商品购物车
		GoodsRouter.PATCH("/:id", order.ShopUpdate)  //修改商品购物车
	}
}
