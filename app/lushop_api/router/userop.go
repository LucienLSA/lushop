package router

import (
	"lushopapi/api/userop"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserOpRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("op")
	AddressRouter := ApiGroup.Group("address").Use(middlewares.JWTAuth())
	{
		AddressRouter.GET("", userop.ARessList)          //获取当前用户的收货地址
		AddressRouter.DELETE("/:id", userop.ARessDelete) //删除收货地址
		AddressRouter.POST("", userop.ARessCreate)       //新建收货地址
		AddressRouter.PUT("/:id", userop.ARessUpdate)    //修改收货地址
	}
	MessageRouter := ApiGroup.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", userop.MsgList)    //获取当前用户的消息
		MessageRouter.POST("", userop.MsgCreate) //新建消息
	}
	UserFavRouter := ApiGroup.Group("userfav").Use(middlewares.JWTAuth())
	{
		UserFavRouter.DELETE("/:id", userop.FavDelete) //删除收藏记录
		UserFavRouter.GET("/:id", userop.FavDetail)    //获取收藏记录
		UserFavRouter.POST("", userop.FavCreate)       //新建收藏记录
		UserFavRouter.GET("", userop.FavList)          //获取当前用户的收藏
	}
}
