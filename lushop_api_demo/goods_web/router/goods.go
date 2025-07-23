package router

import (
	"goodsweb/api/goods"
	"goodsweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("/list", goods.List)
		GoodsRouter.POST("/new", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)                                                             //获取商品详情
		GoodsRouter.DELETE("/delete/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) // 删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)
		GoodsRouter.PATCH("/update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
		GoodsRouter.PUT("/update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
	}
}
