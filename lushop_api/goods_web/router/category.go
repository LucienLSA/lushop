package router

import (
	"lushopapi/goods_web/api/category"
	"lushopapi/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("category")
	{
		CategoryRouter.GET("/list", category.List)
		CategoryRouter.POST("/new", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.New)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.DELETE("/delete/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Delete)
		CategoryRouter.PUT("/update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), category.Update)
	}
}
