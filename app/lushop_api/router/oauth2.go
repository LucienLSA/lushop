package router

import (
	"lushopapi/api/oauth2"

	"github.com/gin-gonic/gin"
)

func InitOAuth2Router(Router *gin.Engine) {
	oauth2Router := Router.Group("/oauth2")
	{
		oauth2Router.GET("/authorize", oauth2.AuthorizeHandler) // 授权
		oauth2Router.POST("/authorize", oauth2.AuthorizeHandler)
		oauth2Router.GET("/login", oauth2.LoginHandler) // 登录
		oauth2Router.POST("/login", oauth2.LoginHandler)
		oauth2Router.GET("/logout", oauth2.LogoutHandler) // 登出
		oauth2Router.POST("/token", oauth2.TokenHandler)  // 获取token或者刷新token
		oauth2Router.GET("/verify", oauth2.VerifyHandler) // 验证token
	}
	Router.NoRoute(oauth2.NotFoundHandler)
}
