package middlewares

import (
	"context"
	"errors"
	"lushopapi/global"
	"lushopapi/utils/jwtClaims"
	"strconv"
	"strings"
	"time"

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
				// 标准化 401 过期响应，便于前端拦截并调用刷新接口
				c.Header("WWW-Authenticate", "Bearer realm=\"api\", error=\"invalid_token\", error_description=\"access token expired\"")
				c.Header("X-Token-Expired", "true")
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":   40101,
					"msg":    "登录授权已过期，请刷新token",
					"action": "refresh_token",
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

		// 检查token是否在黑名单中（仅使用裸token字符串）
		exist, _ := global.RedisClient.Exists(c, "jwt_blacklist:"+tokenString).Result()
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
	// 1) 黑名单检查refresh是否失效
	if exist, _ := global.RedisClient.Exists(context.Background(), "jwt_blacklist:"+rtoken).Result(); exist == 1 {
		return "", "", errors.New("refresh token已失效")
	}

	// 2) 解析并校验 refresh token 如果无效或者过期返回
	rclaim, rErr := j.ParseToken(rtoken)
	if rErr != nil {
		return "", "", errors.New("refresh token无效: " + rErr.Error())
	}
	if rclaim.StandardClaims == nil || time.Now().Unix() > rclaim.StandardClaims.ExpiresAt {
		return "", "", TokenExpired
	}

	// 3) 校验 refresh token 与 Redis 中当前有效值一致（绑定 user + device）
	redisKey := "user:" + strconv.Itoa(int(rclaim.ID)) + ":device:" + rclaim.DeviceID
	storedRToken, getErr := global.RedisClient.Get(context.Background(), redisKey).Result()
	if getErr != nil || storedRToken == "" || storedRToken != rtoken {
		return "", "", errors.New("refresh token不匹配或已失效")
	}

	// 4) 安全解析旧 access token（允许过期，仅校验签名，获取claims）
	aclaims := &jwtClaims.CustomClaims{}
	atok, aErr := jwt.ParseWithClaims(atoken, aclaims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if aErr != nil {
		if ve, ok := aErr.(*jwt.ValidationError); ok {
			// 仅放过过期错误，其它错误一律拒绝
			if ve.Errors&jwt.ValidationErrorExpired == 0 {
				return "", "", errors.New("access token验证失败: " + aErr.Error())
			}
		} else {
			return "", "", errors.New("access token验证失败: " + aErr.Error())
		}
	}
	if atok == nil {
		return "", "", errors.New("access token解析失败")
	}
	// 用户一致性校验
	if aclaims.ID != rclaim.ID {
		return "", "", errors.New("access/refresh token用户不一致，拒绝刷新")
	}

	// 5) 生成新的 access token
	now := time.Now().Unix()
	if aclaims.StandardClaims == nil {
		aclaims.StandardClaims = &jwt.StandardClaims{}
	}
	aclaims.StandardClaims.NotBefore = now
	aclaims.StandardClaims.ExpiresAt = now + global.ServerConfig.JwtInfo.AccessExpireTime
	aclaims.StandardClaims.Issuer = global.ServerConfig.JwtInfo.Key
	newAtoken, err = j.CreateToken(*aclaims)
	if err != nil {
		return "", "", err
	}

	// 6) 旋转 refresh token
	oldRTExp := rclaim.StandardClaims.ExpiresAt
	if rclaim.StandardClaims == nil {
		rclaim.StandardClaims = &jwt.StandardClaims{}
	}
	rclaim.StandardClaims.NotBefore = now
	rclaim.StandardClaims.ExpiresAt = now + global.ServerConfig.JwtInfo.RefreshExpireTime
	rclaim.StandardClaims.Issuer = global.ServerConfig.JwtInfo.Key
	newRtoken, err = j.CreateToken(*rclaim)
	if err != nil {
		return "", "", err
	}
	if setErr := global.RedisClient.Set(context.Background(), redisKey, newRtoken, time.Duration(global.ServerConfig.JwtInfo.RefreshExpireTime)*time.Second).Err(); setErr != nil {
		return "", "", setErr
	}
	// 将旧RT加入黑名单（使用旧的过期时间）
	tl := time.Until(time.Unix(oldRTExp, 0))
	if tl < 0 {
		tl = 0
	}
	_ = global.RedisClient.Set(context.Background(), "jwt_blacklist:"+rtoken, 1, tl).Err()

	return newAtoken, newRtoken, nil
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
