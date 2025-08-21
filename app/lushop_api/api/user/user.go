package user

import (
	"context"
	"fmt"
	"lushopapi/api/base"
	"lushopapi/utils/device"
	"lushopapi/utils/jwtClaims"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"lushopapi/forms"
	"lushopapi/global"
	"lushopapi/global/response"
	"lushopapi/middlewares"
	v2userproto "lushopapi/proto/user"
)

// 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 1. 权限验证：仅管理员可操作
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	adminUser := claims.(*jwtClaims.CustomClaims)
	fmt.Println(adminUser)
	if adminUser.AuthorityId != 2 { // 假设2为管理员权限ID
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无管理员权限"})
		return
	}
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &v2userproto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

}

// 密码登录
func PassWorldLogin(c *gin.Context) {
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		base.HandleValidatorError(c, err)
		return
	}
	// // 校验验证码时从 Redis 取出
	// realAns, err := global.RedisClient.Get(c, "captcha:"+passwordLoginForm.CaptchaId).Result()
	// if err != nil {
	// 	// 过期或不存在
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "验证码过期或不存在",
	// 	})
	// 	return
	// }
	// if realAns != passwordLoginForm.CaptchaAns {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "验证码错误",
	// 	})
	// 	return
	// }
	// // 校验通过后删除，防止重放
	// global.RedisClient.Del(c, "captcha:"+passwordLoginForm.CaptchaId)

	// 直接通过Captcha自带的store进行检验
	// if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.CaptchaAns, true) {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"captcha": "验证码错误",
	// 	})
	// 	return
	// }
	fmt.Printf("输入的手机号为:%s", passwordLoginForm.Mobile)
	// fmt.Printf("上下文保存的手机号为:%v", c)
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.WithValue(context.Background(), "ginContext", c), &v2userproto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				//zap.S().Error(err)
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
					"code":   strconv.Itoa(int(e.Code())),
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &v2userproto.PasswordCheckInfo{
			PassWord:          passwordLoginForm.PassWord, // 前端用户传入的密码
			EncryptedPassWord: rsp.PassWord,               // 数据库中查询到的用户设置的密码
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "密码错误",
			})
		} else {
			if passRsp.Success {
				// 生成access_token
				j := middlewares.NewJWT()
				aclaims := jwtClaims.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: &jwt.StandardClaims{
						NotBefore: time.Now().Unix(),                                                //签名的生效时间
						ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.AccessExpireTime, //30天过期
						Issuer:    global.ServerConfig.JwtInfo.Key,
					},
				}
				// 生成access_token
				AccessToken, err := j.CreateToken(aclaims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成access_token失败",
					})
					return
				}
				// 将生成access_token存储到Redis，键为用户ID，过期时间与生成access_token有效期一致
				err = global.RedisClient.Set(c, "user_token:"+strconv.Itoa(int(rsp.Id)),
					AccessToken, time.Duration(global.ServerConfig.JwtInfo.AccessExpireTime)*time.Second).Err()
				if err != nil {
					zap.S().Errorf("存储token到Redis失败: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录异常"})
					return
				}
				// 生成refresh_token 不需要保存任何用户信息
				rclaims := jwtClaims.CustomClaims{
					ID:       uint(rsp.Id),                     // 绑定用户ID到refresh token
					DeviceID: device.GetDeviceFingerprint(c), // 添加设备指纹
					StandardClaims: &jwt.StandardClaims{
						NotBefore: time.Now().Unix(),                                                 //签名的生效时间
						ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.RefreshExpireTime, //30天过期
						Issuer:    global.ServerConfig.JwtInfo.Key,
					},
				}
				RefreshToken, err := j.CreateToken(rclaims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成refresh_token失败",
					})
					return
				}
				// 将refreshtoken存储到redis中，用刷新accesstoken时验证设备
				redisKey := fmt.Sprintf("user:%d:device:%s", rsp.Id, rclaims.DeviceID)
				err = global.RedisClient.Set(c, redisKey, RefreshToken, time.Duration(global.ServerConfig.JwtInfo.RefreshExpireTime)*time.Second).Err()
				if err != nil {
					zap.S().Errorf("存储refresh token失败: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"msg": "登录异常"})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":            rsp.Id,
					"nick_name":     rsp.NickName,
					"access_token":  AccessToken,
					"refresh_token": RefreshToken,
					"role":          rsp.Role,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登陆失败",
				})
			}
		}
	}
}

// 登出
func Logout(c *gin.Context) {
	// 获取当前用户token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "未携带token"})
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Authorization格式错误"})
		return
	}
	token := parts[1]

	// 获取当前用户claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户：%v", currentUser)
	// 解析token获取过期时间
	// 计算剩余有效期作为黑名单过期时间
	expire := time.Unix(currentUser.ExpiresAt, 0).Sub(time.Now())

	// 将jwt加入黑名单
	global.RedisClient.Set(c, "jwt_blacklist:"+token, 1, expire)
	c.JSON(http.StatusOK, gin.H{"msg": "注销成功"})
}

// RefreshTokenHandler 刷新accessToken
func RefreshToken(c *gin.Context) {
	refreshFroms := forms.RefreshTokenForm{}
	if err := c.ShouldBind(&refreshFroms); err != nil {
		base.HandleValidatorError(c, err)
		return
	}
	// 从 Authorization: Bearer <token> 中提取旧 access token
	authHeader := c.Request.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"msg": "Authorization格式错误",
		})
		c.Abort()
		return
	}
	oldAccessToken := parts[1]
	j := middlewares.NewJWT()
	// 解析refresh token获取设备指纹
	refreshClaims, err := j.ParseToken(refreshFroms.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "refresh token无效"})
		return
	}
	// 验证当前设备指纹是否匹配
	currentDeviceID := device.GetDeviceFingerprint(c)
	if refreshClaims.DeviceID != currentDeviceID {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "设备验证失败，拒绝刷新token"})
		return
	}
	//调用刷新token方法，生成新的AccessToken, RefreshToken
	AccessToken, RefreshToken, err := j.RefreshToken(oldAccessToken, refreshFroms.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": "刷新双Token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  AccessToken,
		"refresh_token": RefreshToken, // 临时做测试返回
	})
	// HttpOnly Cookie传输令牌
	c.SetCookie("refresh_token", RefreshToken, int(global.ServerConfig.JwtInfo.RefreshExpireTime), "/", "localhost", false, true)
	c.SetCookie("access_token", AccessToken, int(global.ServerConfig.JwtInfo.AccessExpireTime), "/", "localhost", false, true)
}

func Register(c *gin.Context) {
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		base.HandleValidatorError(c, err)
		return
	}
	//验证码
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// })
	value, err := global.RedisClient.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码为空",
		})
		return
	} else {
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}
	}

	//生成grpc的client并调用接口
	user, err := global.UserSrvClient.CreateUser(context.WithValue(context.Background(), "ginContext", c), &v2userproto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 新建 【新建用户失败】失败：%s", err.Error())
		base.HandleGrpcErrorToHttp(err, c)
		return
	}
	// 用户注册完成后，需要自行登录
	// //生成token
	// j := middlewares.NewJWT()
	// claims := jwtClaims.CustomClaims{
	// 	ID:          uint(user.Id),
	// 	NickName:    user.NickName,
	// 	AuthorityId: uint(user.Role),
	// 	StandardClaims: &jwt.StandardClaims{
	// 		NotBefore: time.Now().Unix(),                                          //签名的生效时间
	// 		ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime, //30天过期
	// 		Issuer:    "Lushop",
	// 	},
	// }
	// token, err := j.CreateToken(claims)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"msg": "生成token失败",
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"nick_name": user.NickName,
		// "token":       token,
		// "expiresd_at": (time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime) * 1000,
	})
}

func GetUserDetailOAuth(c *gin.Context) {
	// claims, _ := c.Get("claims")
	// currentUser := claims.(*jwtClaims.CustomClaims)
	// zap.S().Infof("访问用户：%d", currentUser)

	userID, exists := c.Get("client_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	mobile := userID.(string)
	// 这里才去找用户的实际id
	user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, c)
		return
	}
	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &v2userproto.IdRequest{
		// Id: int32(currentUser.ID),
		Id: int32(user.Id),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":     rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}

// 获取当前用户详情
func GetUserDetail(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户：%v", currentUser)

	// 这里才去找用户的实际id
	// user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
	// 	Mobile: mobile,
	// })
	// if err != nil {
	// 	base.HandleGrpcErrorToHttp(err, c)
	// 	return
	// }
	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &v2userproto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":     rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}

// 更新用户
func UpdateUserOAuth(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	// claims, _ := ctx.Get("claims")
	// currentUser := claims.(*jwtClaims.CustomClaims)
	// zap.S().Infof("访问用户：%d", currentUser)
	userID, exists := ctx.Get("client_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	mobile := userID.(string)
	// 这里才去找用户的实际id
	user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	fmt.Println(user)
	// rsp, err := global.UserSrvClient.GetUserById(context.Background(), &v2userproto.IdRequest{
	// 	// Id: int32(currentUser.ID),
	// 	Id: int32(user.Id),
	// })
	// if err != nil {
	// 	base.HandleGrpcErrorToHttp(err, ctx)
	// 	return
	// }
	//将前端传递过来的日期格式转换成int类型
	loc, _ := time.LoadLocation("Local")                                            // L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc) //必须是2006-01-02
	_, err = global.UserSrvClient.UpdateUser(context.Background(), &v2userproto.UpdateUserInfo{
		// Id:       int32(currentUser.ID),
		Id:       user.Id,
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新用户信息成功",
	})
}

func UpdateUserInfo(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser)

	// // 这里才去找用户的实际id
	// user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
	// 	Mobile: mobile,
	// })
	// if err != nil {
	// 	base.HandleGrpcErrorToHttp(err, ctx)
	// 	return
	// }

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &v2userproto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//将前端传递过来的日期格式转换成int类型
	loc, _ := time.LoadLocation("Local")                                            // L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc) //必须是2006-01-02
	_, err = global.UserSrvClient.UpdateUser(context.Background(), &v2userproto.UpdateUserInfo{
		Id: rsp.Id,
		// Id:       user.Id,
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新用户信息成功",
	})
}

// UpdatePassword 修改密码接口
func UpdatePassword(c *gin.Context) {
	updatePwdForm := forms.UpdatePasswordForm{}
	if err := c.ShouldBind(&updatePwdForm); err != nil {
		base.HandleValidatorError(c, err)
		return
	}

	// 获取当前用户ID
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	currentUser := claims.(*jwtClaims.CustomClaims)
	userID := currentUser.ID

	// 1. 验证旧密码
	user, err := global.UserSrvClient.GetUserById(context.Background(), &v2userproto.IdRequest{
		Id: int32(userID),
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, c)
		return
	}

	// 2. 检查旧密码是否正确
	passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &v2userproto.PasswordCheckInfo{
		PassWord:          updatePwdForm.OldPassword,
		EncryptedPassWord: user.PassWord,
	})
	if pasErr != nil || !passRsp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
		return
	}

	// 3. 更新新密码
	_, err = global.UserSrvClient.UpdatePassword(context.Background(), &v2userproto.UpdatePasswordInfo{
		Id:          int32(userID),
		NewPassword: updatePwdForm.NewPassword,
	})
	if err != nil {
		base.HandleGrpcErrorToHttp(err, c)
		return
	}

	// 4. Token失效处理
	// 将当前access_token加入黑名单并删除存储
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			authHeader = parts[1]
		}
		j := middlewares.NewJWT()
		claims, err := j.ParseToken(authHeader)
		if err == nil {
			expire := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
			global.RedisClient.Set(c, "jwt_blacklist:"+authHeader, 1, expire)
		}
	}
	// 删除用户access_token
	err = global.RedisClient.Del(c, "user_token:"+strconv.Itoa(int(userID))).Err()
	if err != nil {
		zap.S().Errorf("删除access_token失败: %v", err)
	}
	// 将refresh_token加入到黑名单
	// 删除所有设备refresh_token
	iter := global.RedisClient.Scan(c, 0, "user:"+strconv.Itoa(int(userID))+":device:*", 0).Iterator()
	for iter.Next(c) {
		key := iter.Val()
		token, _ := global.RedisClient.Get(c, key).Result()
		if token != "" {
			j := middlewares.NewJWT()
			claims, err := j.ParseToken(token)
			if err == nil {
				expire := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
				global.RedisClient.Set(c, "jwt_blacklist:"+token, 1, expire)
			}
		}
		global.RedisClient.Del(c, key)
	}
	if err := iter.Err(); err != nil {
		zap.S().Errorf("删除refresh_token失败: %v", err)
	}
	//强制用户重新登录
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusUnauthorized, gin.H{
		"error":   "password_changed",
		"message": "密码已修改，请使用新密码重新登录",
	})
}

// ForceLogout 管理员强制用户下线
// @Summary 强制用户下线
// @Tags 管理员接口
// @Accept json
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} gin.H{"msg":"强制下线成功"}
// @Failure 400 {object} gin.H{"error":"参数错误"}
// @Failure 401 {object} gin.H{"error":"未授权"}
// @Failure 403 {object} gin.H{"error":"无权限"}
// @Router /admin/user/{userId}/logout [post]
func ForceLogout(c *gin.Context) {
	// 1. 权限验证：仅管理员可操作
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	adminUser := claims.(*jwtClaims.CustomClaims)
	if adminUser.AuthorityId != 2 { // 假设2为管理员权限ID
		c.JSON(http.StatusForbidden, gin.H{"error": "无管理员权限"})
		return
	}

	// 2. 获取目标用户ID
	userIdStr := c.Query("userId")

	// 3. 使目标用户所有Token失效
	ctx := context.Background()

	// 3.1 处理access_token
	accessTokenKey := "user_token:" + userIdStr
	// 获取当前access_token并加入黑名单
	token, err := global.RedisClient.Get(ctx, accessTokenKey).Result()
	if err == nil && token != "" {
		j := middlewares.NewJWT()
		claims, err := j.ParseToken(token)
		if err == nil {
			expire := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
			global.RedisClient.Set(ctx, "jwt_blacklist:"+token, 1, expire)
		}
	}
	fmt.Println(token)
	// 删除access_token存储
	global.RedisClient.Del(ctx, accessTokenKey)

	// 处理refresh_token
	iter := global.RedisClient.Scan(ctx, 0, "user:"+userIdStr+":device:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		token, _ := global.RedisClient.Get(ctx, key).Result()
		if token != "" {
			j := middlewares.NewJWT()
			claims, err := j.ParseToken(token)
			// 将refresh_token加入黑名单
			if err == nil {
				expire := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
				global.RedisClient.Set(ctx, "jwt_blacklist:"+token, 1, expire)
			}
		}
		// 删除refresh_token redis存储
		global.RedisClient.Del(ctx, key)
	}
	if err := iter.Err(); err != nil {
		zap.S().Errorf("强制下线扫描refresh_token失败: %v", err)
	}

	zap.S().Infof("管理员[%d]强制用户[%s]下线成功", adminUser.ID, userIdStr)
	c.JSON(http.StatusOK, gin.H{"msg": "强制下线成功"})
}
