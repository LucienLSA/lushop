package router

import (
	v2user "lushopapi/api/user"
	"lushopapi/global"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("u")
	BaseRouter := ApiGroup.Group("num")
	{
		BaseRouter.GET("captcha", v2user.GetCaptcha)
		BaseRouter.POST("send_sms", v2user.SendSmsAli)
	}
}
