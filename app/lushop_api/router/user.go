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
		UserRouer.POST("pwd_login", user.PassWorldLogin)
		UserRouer.POST("register", user.Register)
		UserRouer.GET("refresh", user.RefreshToken)
		UserRouer.GET("detail", middlewares.JWTAuth(), user.GetUserDetail)
		UserRouer.PATCH("update", middlewares.JWTAuth(), user.UpdateUser)
		UserRouer.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), user.GetUserList)
	}

	// OAuth2方案
	UserRouerV2 := Router.Group("user2")
	{
		UserRouerV2.POST("pwd_login", user.PassWorldLogin)
		UserRouerV2.POST("register", user.Register)
		UserRouerV2.GET("list", middlewares.AuthMiddleware(), user.GetUserList)
		UserRouerV2.GET("detail", middlewares.AuthMiddleware(), user.GetUserDetail)
		UserRouerV2.PATCH("update", middlewares.AuthMiddleware(), user.UpdateUser)
	}
}
