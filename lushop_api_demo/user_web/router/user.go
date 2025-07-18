package router

import (
	"userweb/api"
	"userweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
		UserRouter.POST("refresh_token", middlewares.JWTAuth(), api.RefreshToken)
		UserRouter.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
}
