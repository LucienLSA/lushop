package router

import (
	"lushopapi/api/oss"

	"github.com/gin-gonic/gin"
)

func InitOssRouter(Router *gin.Engine) {
	OssRouter := Router.Group("oss")
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", oss.Token)
		OssRouter.POST("/callback", oss.HandlerRequest)
	}
}
