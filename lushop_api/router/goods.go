package router

import (
	"lushopapi/api/goods"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.Engine) {
	ApiGroup := Router.Group("/g")
	ApiGroup = ApiGroup.Group(global.ServerConfig.Version)
	GoodsRouter := ApiGroup.Group("/goods")
	{
		// GoodsRouter.GET("", goods.GoodsList)                                                                 //商品列表(es)
		GoodsRouter.GET("", goods.GoodsListES)                                                               //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.GoodsCreate)            //新建商品
		GoodsRouter.GET("/:id", goods.GoodsDetail)                                                           //商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.GoodsDelete)      //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                                         //获取库存
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.GoodsUpdate)         //更新库存
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.GoodsUpdateStatus) //商品状态

	}

	// BannerRouter := Router.Group("banner").Use(middlewares.Trace())
	BannerRouter := ApiGroup.Group("banners")
	{
		BannerRouter.GET("", goods.BannerList)                                                            // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BannerDelete) // 删除轮播图
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BannerCreate)       //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BannerUpdate)    //修改轮播图信息
	}

	BrandRouter := ApiGroup.Group("brands")
	{
		BrandRouter.GET("", goods.BrandList)                                                            // 品牌列表页
		BrandRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BrandDelete) // 删除品牌
		BrandRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BrandCreate)       //新建品牌
		BrandRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.BrandUpdate)    //修改品牌信息
	}

	CategoryBrandRouter := ApiGroup.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", goods.CateBrandList)                                                            // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateBrandDelete) // 删除类别品牌
		CategoryBrandRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateBrandCreate)       //新建类别品牌
		CategoryBrandRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateBrandUpdate)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", goods.CateGetBrandList)                                                     //获取分类的品牌
	}

	CategoryRouter := ApiGroup.Group("categorys")
	{
		CategoryRouter.GET("", goods.CateList)                                                            // 商品类别列表页
		CategoryRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateDelete) // 删除分类
		CategoryRouter.GET("/:id", goods.CateDetail)                                                      // 获取分类详情
		CategoryRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateCreate)       //新建分类
		CategoryRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.CateUpdate)    //修改分类信息
	}
}
