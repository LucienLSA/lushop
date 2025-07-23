package router

import (
	"lushopapi/api/userop"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserOpRouter(Router *gin.Engine) {
	AddressRouter := Router.Group("address")
	{
		AddressRouter.GET("", middlewares.JWTAuth(), userop.ARessList)          //获取当前用户的收货地址
		AddressRouter.DELETE("/:id", middlewares.JWTAuth(), userop.ARessDelete) //删除收货地址
		AddressRouter.POST("", middlewares.JWTAuth(), userop.ARessCreate)       //新建收货地址
		AddressRouter.PUT("/:id", middlewares.JWTAuth(), userop.ARessUpdate)    //修改收货地址
	}
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", userop.MsgList)    //获取当前用户的消息
		MessageRouter.POST("", userop.MsgCreate) //新建消息
	}
	UserFavRouter := Router.Group("userfav")
	{
		UserFavRouter.DELETE("/:id", middlewares.JWTAuth(), userop.FavDelete) //删除收藏记录
		UserFavRouter.GET("/:id", middlewares.JWTAuth(), userop.FavDetail)    //获取收藏记录
		UserFavRouter.POST("", middlewares.JWTAuth(), userop.FavCreate)       //新建收藏记录
		UserFavRouter.GET("", middlewares.JWTAuth(), userop.FavList)          //获取当前用户的收藏
	}
}
