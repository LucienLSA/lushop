package middlewares

import (
	"errors"
	"lushopapi/global"
	"lushopapi/utils/jwtClaims"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 Authorization: Bearer <token>
		// 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			zap.S().Errorf("Authorization header is empty")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请登录",
			})
			c.Abort()
			return
		}

		// 添加调试信息
		// zap.S().Infof("收到Authorization Header: %s", authHeader)

		// 按空格分割，检查Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		// zap.S().Infof("分割后的部分: %v, 长度: %d", parts, len(parts))

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			// zap.S().Errorf("Authorization header format error: %s", authHeader)
			// zap.S().Errorf("期望格式: Bearer <token>, 实际格式: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization格式错误, 应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		if tokenString == "" {
			zap.S().Errorf("Token string is empty")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token不能为空",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		// 解析token
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			switch err {
			case TokenExpired:
				zap.S().Errorf("Token expired: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "登录授权已过期",
				})
			case TokenMalformed:
				zap.S().Errorf("Token malformed: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "Token格式错误",
				})
			case TokenNotValidYet:
				zap.S().Errorf("Token not valid yet: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "Token尚未生效",
				})
			default:
				zap.S().Errorf("Token validation error: %s", err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "Token验证失败",
				})
			}
			c.Abort()
			return
		}

		// 检查token是否在黑名单中
		exist, _ := global.RedisClient.Exists(c, "jwt_blacklist:"+authHeader).Result()
		if exist == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token已失效",
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JwtInfo.Secret),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims jwtClaims.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*jwtClaims.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
		// 如果不是ValidationError，直接返回TokenInvalid
		return nil, TokenInvalid
	}
	if token != nil {
		if claims, ok := token.Claims.(*jwtClaims.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// RefreshToken 通过 refresh token 刷新 atoken
func (j *JWT) RefreshToken(atoken, rtoken string) (newAtoken, newRtoken string, err error) {
	// rtoken 无效直接返回
	var rclaim *jwtClaims.CustomClaims
	if rclaim, err = j.ParseToken(rtoken); err != nil {
		return
	}
	// 从旧access token 中解析出claims数据
	var aclaim *jwtClaims.CustomClaims
	aclaim, aErr := j.ParseToken(atoken)
	if aErr != nil {
		// 判断错误是不是因为access token 正常过期导致的
		if v, ok := aErr.(*jwt.ValidationError); ok {
			if v.Errors == jwt.ValidationErrorExpired {
				// 刷新生成新的access_token
				NewAccess, _ := j.CreateToken(*aclaim)
				NewRefresh, _ := j.CreateToken(*rclaim)
				return NewAccess, NewRefresh, nil
			}
		}
		// 其他错误直接返回
		err = aErr
		return
	}
	// access token 没有过期，不需要刷新
	zap.S().Info("access token 未过期，无需刷新")
	return
}

// 更新token
// func (j *JWT) RefreshToken(tokenString string) (string, error) {
// 	jwt.TimeFunc = func() time.Time {
// 		return time.Unix(0, 0)
// 	}
// 	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return j.SigningKey, nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if claims, ok := token.Claims.(*jwtClaims.CustomClaims); ok && token.Valid {
// 		jwt.TimeFunc = time.Now
// 		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
// 		return j.CreateToken(*claims)
// 	}
// 	return "", TokenInvalid
// }
