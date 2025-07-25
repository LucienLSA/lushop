package router

import (
	v2user "lushopapi/api/user"
	"lushopapi/global"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.Engine) {
	ApiGroup := Router.Group("u")
	ApiGroup = ApiGroup.Group(global.ServerConfig.Version)
	BaseRouter := ApiGroup.Group("base")
	{
		BaseRouter.GET("captcha", v2user.GetCaptcha)
		BaseRouter.POST("send_sms", v2user.SendSmsAli)
	}
}
