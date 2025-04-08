package router

import (
	"useropweb/api/user_fav"
	"useropweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs")
	{
		UserFavRouter.DELETE("/delete/:id", middlewares.JWTAuth(), user_fav.Delete) // 删除收藏记录
		UserFavRouter.GET("/detail/:id", middlewares.JWTAuth(), user_fav.Detail)    // 获取收藏记录
		UserFavRouter.POST("/new", middlewares.JWTAuth(), user_fav.New)             //新建收藏记录
		UserFavRouter.GET("/list", middlewares.JWTAuth(), user_fav.List)            //获取当前用户的收藏
	}
}
