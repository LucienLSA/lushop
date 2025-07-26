package router

import (
	"lushopapi/api/oss"
	"lushopapi/global"

	"github.com/gin-gonic/gin"
)

func InitOssRouter(Router *gin.Engine) {
	ApiGroup := Router.Group("g")
	ApiGroup = ApiGroup.Group(global.ServerConfig.Version)
	OssRouter := ApiGroup.Group("oss")
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", oss.Token)              // 获取上传策略令牌
		OssRouter.POST("callback", oss.HandlerRequest) //回调验证
	}
}
