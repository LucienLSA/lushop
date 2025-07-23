package router

import (
	"orderweb/api/shopCart"
	"orderweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcart").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("/list", shopCart.List)            //购物车列表
		GoodsRouter.DELETE("/delete/:id", shopCart.Delete) //删除条目
		GoodsRouter.POST("/new", shopCart.New)             //添加商品到购物车
		GoodsRouter.PATCH("/update/:id", shopCart.Update)  //修改条目
	}
}
