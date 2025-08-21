package router

import (
	"lushopapi/api/user"
	"lushopapi/global"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.Engine) {
	ApiGroup := Router.Group(global.ServerConfig.Version)
	// ApiGroup = ApiGroup.Group("u")
	UserRouer := ApiGroup.Group("user")
	zap.S().Info("配置用户相关的url")
	// 双token方案
	{
		// UserRouer.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), user.GetUserList)
		UserRouer.POST("pwd_login", user.PassWorldLogin)                                             // 用户账号密码登录
		UserRouer.POST("register", user.Register)                                                    // 用户注册
		UserRouer.POST("refresh", user.RefreshToken)                                                 // 用户刷新token
		UserRouer.PUT("update_pwd", middlewares.JWTAuth(), user.UpdatePassword)                      // 用户密码灯芯
		UserRouer.GET("detail", middlewares.JWTAuth(), user.GetUserDetail)                           // 用户详情
		UserRouer.PATCH("update_info", middlewares.JWTAuth(), user.UpdateUserInfo)                   // 用户信息更新
		UserRouer.DELETE("logout", middlewares.JWTAuth(), user.Logout)                               // 用户登出
		UserRouer.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), user.GetUserList)    // 管理员获取用户列表
		UserRouer.DELETE("kick", middlewares.JWTAuth(), middlewares.IsAdminAuth(), user.ForceLogout) // 管理员踢人出系统
	}

	// OAuth2方案
	UserRouerV2 := Router.Group("user2")
	{
		UserRouerV2.POST("pwd_login", user.PassWorldLogin)
		UserRouerV2.POST("register", user.Register)
		UserRouerV2.GET("list", middlewares.AuthMiddleware(), user.GetUserList)
		UserRouerV2.GET("detail", middlewares.AuthMiddleware(), user.GetUserDetail)
		UserRouerV2.PATCH("update", middlewares.AuthMiddleware(), user.UpdateUserInfo)
	}

	// 验证码路由
	CaptchaRouter := ApiGroup.Group("captcha")
	{
		CaptchaRouter.GET("picture/get", user.GetCaptchaV2) // 获取图形验证码
		// CaptchaRouter.GET("picture/refresh", user.RefreshCaptcha) // 刷新图形验证码
		CaptchaRouter.GET("picture/verify", user.VerifyCaptcha) // 验证图形验证码

		CaptchaRouter.POST("sms/send", user.SendSmsAli) // 发送手机验证码
		// CaptchaRouter.GET("sms/refresh", user.RefreshSmsAli) // 刷新手机验证码
		CaptchaRouter.GET("sms/verify", user.VerifySmsAli) // 刷新手机验证码
	}
}
