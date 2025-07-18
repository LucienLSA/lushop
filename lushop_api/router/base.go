package router

import (
	v2user "lushopapi/api/user"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.Engine) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", v2user.GetCaptcha)
		BaseRouter.POST("send_sms", v2user.SendSmsAli)
	}
}
