package router

import (
	"lushopapi/order_web/api/shopcart"
	"lushopapi/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcart").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("/list", shopcart.List)            //购物车列表
		GoodsRouter.DELETE("/delete/:id", shopcart.Delete) //删除条目
		GoodsRouter.POST("/new", shopcart.New)             //添加商品到购物车
		GoodsRouter.PATCH("/update/:id", shopcart.Update)  //修改条目
	}
}
