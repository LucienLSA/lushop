package router

import (
	"lushopapi/goods_web/api/banner"
	"lushopapi/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("/list", banner.List)                                                              // 轮播图列表页
		BannerRouter.DELETE("/delete/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.Delete) // 删除轮播图
		BannerRouter.POST("/new", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.New)             //新建轮播图
		BannerRouter.PUT("/update/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banner.Update)    //修改轮播图信息
	}
}
