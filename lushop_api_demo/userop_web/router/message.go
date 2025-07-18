package router

import (
	"useropweb/api/message"
	"useropweb/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(Router *gin.RouterGroup) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("/list", message.List)
		MessageRouter.POST("/new", message.New)
	}
}
