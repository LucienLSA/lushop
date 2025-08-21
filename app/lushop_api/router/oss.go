package router

import (
	"lushopapi/api/oss"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOssRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("g")
	OssRouter := ApiGroup.Group("oss").Use(middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", oss.Token)              // 获取上传策略令牌
		OssRouter.POST("callback", oss.HandlerRequest) //回调验证
		OssRouter.POST("upload", oss.PostPicture)
	}
}
