package router

import (
	"goodsweb/api/brand"

	"github.com/gin-gonic/gin"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands")
	{
		BrandRouter.GET("/list", brand.BrandList)            // 品牌列表页
		BrandRouter.DELETE("/delete/:id", brand.DeleteBrand) // 删除品牌
		BrandRouter.POST("/new", brand.NewBrand)             //新建品牌
		BrandRouter.PUT("/update/:id", brand.UpdateBrand)    //修改品牌信息
	}
	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("/list", brand.CategoryBrandList)            // 类别品牌列表页
		CategoryBrandRouter.DELETE("/delete/:id", brand.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("/new", brand.NewCategoryBrand)             //新建类别品牌
		CategoryBrandRouter.PUT("/update:id", brand.UpdateCategoryBrand)     //修改类别品牌
		CategoryBrandRouter.GET("/list/:id", brand.GetCategoryBrandList)     //获取分类的品牌
	}
}
