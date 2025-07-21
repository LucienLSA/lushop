package user

import (
	"context"
	"lushopapi/api/base"
	"lushopapi/utils/jwtClaims"
	"net/http"
	"strconv"
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
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.CaptchaAns, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &v2userproto.MobileRequest{
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
					"mobile": "登录失败-1",
					"code":   strconv.Itoa(int(e.Code())),
				})
			}
			return
		}
	} else {
		//只是查询到用户了而已，并没有检查密码
		if passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &v2userproto.PasswordCheckInfo{
			PassWord:          passwordLoginForm.PassWord,
			EncryptedPassWord: rsp.PassWord,
		}); pasErr != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"password": "密码错误",
			})
		} else {
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := jwtClaims.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: &jwt.StandardClaims{
						NotBefore: time.Now().Unix(),                                          //签名的生效时间
						ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime, //30天过期
						Issuer:    "lucien",
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":          rsp.Id,
					"nick_name":   rsp.NickName,
					"token":       token,
					"expiresd_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登陆失败",
				})
			}
		}
	}
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

// 注销用户
// func DeleteUser(c *gin.Context) {
// 	claims, _ := c.Get("claims")
// 	currentUser := claims.(*jwtClaims.CustomClaims)
// 	_, err := global.UserSrvClient.DeleteUser(context.Background(), &v2userproto.IdRequest{
// 		Id: int32(currentUser.ID),
// 	})
// 	if err != nil {
// 		base.HandleGrpcErrorToHttp(err, c)
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"msg": "注销成功",
// 	})
// }

func GetUserDetail(c *gin.Context) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser)
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
func UpdateUser(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户：%d", currentUser)

	//将前端传递过来的日期格式转换成int类型
	loc, _ := time.LoadLocation("Local")                                            // L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc) //必须是2006-01-02
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &v2userproto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
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
