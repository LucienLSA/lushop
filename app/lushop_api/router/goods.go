package router

import (
	"lushopapi/api/goods"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("/g")
	GoodsRouter := ApiGroup.Group("/goods").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("/list", goods.GoodsList)                                                   //商品列表
		GoodsRouter.GET("/list/es", goods.GoodsListES)                                              // ES查询商品列表
		GoodsRouter.POST("/create", middlewares.IsAdminAuth(), goods.GoodsCreate)                   //新建商品
		GoodsRouter.GET("/detail/:id", goods.GoodsDetail)                                           //商品详情
		GoodsRouter.DELETE("/delete/:id", middlewares.IsAdminAuth(), goods.GoodsDelete)             //删除商品
		GoodsRouter.GET("/stocks/get/:id/", goods.Stocks)                                           //获取库存
		GoodsRouter.PUT("/stocks/update/:id", middlewares.IsAdminAuth(), goods.GoodsUpdate)         //更新库存
		GoodsRouter.PATCH("/status/update/:id", middlewares.IsAdminAuth(), goods.GoodsUpdateStatus) //更新商品状态

	}

	// BannerRouter := Router.Group("banner").Use(middlewares.Trace())
	BannerRouter := ApiGroup.Group("banners").Use(middlewares.JWTAuth())
	{
		BannerRouter.GET("/list", goods.BannerList) // 轮播图列表页

		BannerRouter.DELETE("/delete/:id", middlewares.IsAdminAuth(), goods.BannerDelete) // 删除轮播图
		BannerRouter.POST("/cteate", middlewares.IsAdminAuth(), goods.BannerCreate)       //新建轮播图
		BannerRouter.PUT("/update/:id", middlewares.IsAdminAuth(), goods.BannerUpdate)    //修改轮播图信息
	}

	BrandRouter := ApiGroup.Group("brands").Use(middlewares.JWTAuth())
	{
		BrandRouter.GET("/list", goods.BrandList)                                       // 品牌列表页
		BrandRouter.DELETE("/delete/:id", middlewares.IsAdminAuth(), goods.BrandDelete) // 删除品牌
		BrandRouter.POST("/create", middlewares.IsAdminAuth(), goods.BrandCreate)       //新建品牌
		BrandRouter.PUT("/update/:id", middlewares.IsAdminAuth(), goods.BrandUpdate)    //修改品牌信息
	}

	CategoryBrandRouter := ApiGroup.Group("categorybrands").Use(middlewares.JWTAuth())
	{
		CategoryBrandRouter.GET("/list", goods.CateBrandList)                                       // 类别品牌列表页
		CategoryBrandRouter.DELETE("/delete/:id", middlewares.IsAdminAuth(), goods.CateBrandDelete) // 删除类别品牌
		CategoryBrandRouter.POST("/create", middlewares.IsAdminAuth(), goods.CateBrandCreate)       //新建类别品牌
		CategoryBrandRouter.PUT("/update/:id", middlewares.IsAdminAuth(), goods.CateBrandUpdate)    //修改类别品牌
		CategoryBrandRouter.GET("/getbrandlist/:id", goods.CateGetBrandList)                        //获取分类的品牌
	}

	CategoryRouter := ApiGroup.Group("categorys").Use(middlewares.JWTAuth())
	{
		CategoryRouter.GET("/list", goods.CateList)                                       // 商品类别列表页
		CategoryRouter.DELETE("/delete/:id", middlewares.IsAdminAuth(), goods.CateDelete) // 删除分类
		CategoryRouter.GET("/detail/:id", goods.CateDetail)                               // 获取分类详情
		CategoryRouter.POST("/create", middlewares.IsAdminAuth(), goods.CateCreate)       //新建分类
		CategoryRouter.PUT("/update/:id", middlewares.IsAdminAuth(), goods.CateUpdate)    //修改分类信息
	}
}
