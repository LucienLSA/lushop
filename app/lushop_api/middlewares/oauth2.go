package middlewares

import (
	"fmt"
	"lushopapi/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 校验 OAuth2 Bearer Token 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 校验token
		ti, err := global.Srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录或token无效",
			})
			c.Abort()
			return
		}
		fmt.Println(ti.GetUserID(), ti.GetClientID(), ti.GetScope())
		// 将用户信息存入 context，便于后续 handler 使用
		c.Set("user_id", ti.GetUserID())
		c.Set("client_id", ti.GetClientID())
		c.Set("scope", ti.GetScope())
		c.Next()
	}
}
