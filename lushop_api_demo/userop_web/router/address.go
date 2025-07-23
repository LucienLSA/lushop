package router

import (
	"useropweb/api/address"
	"useropweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("/list", middlewares.JWTAuth(), address.List)            // 轮播图列表页
		AddressRouter.DELETE("/delete/:id", middlewares.JWTAuth(), address.Delete) // 删除轮播图
		AddressRouter.POST("/new", middlewares.JWTAuth(), address.New)             //新建轮播图
		AddressRouter.PUT("/update/:id", middlewares.JWTAuth(), address.Update)    //修改轮播图信息
	}
}
