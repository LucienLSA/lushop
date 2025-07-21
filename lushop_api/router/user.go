package router

import (
	"lushopapi/api/user"
	"lushopapi/middlewares"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.Engine) {
	UserRouer := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouer.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), user.GetUserList)
		//UserRouer.GET("list", api.GetUserList)
		UserRouer.POST("pwd_login", user.PassWorldLogin)
		UserRouer.POST("register", user.Register)

		UserRouer.GET("detail", middlewares.JWTAuth(), user.GetUserDetail)
		UserRouer.PATCH("update", middlewares.JWTAuth(), user.UpdateUser)
		// UserRouer.POST("refresh", )
	}
}
