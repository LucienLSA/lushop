package api

import (
	"context"
	"userweb/forms"
	"userweb/global"
	"userweb/global/response"
	"userweb/global/types"
	"userweb/middlewares"
	"userweb/proto"
	"userweb/utils/jwtClaims"

	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message,
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func HandlerValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户:%d", currentUser.ID)
	zap.S().Infof("用户身份:%d", currentUser.AuthorityId)
	// 生成grpc的client并调用接口
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "5")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询【用户列表失败】", err)
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]*response.UserResponse, 0)
	for _, value := range rsp.Data {
		// data := make(map[string]interface{})
		user := &response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(types.Uint64ToTime(value.BirthDay)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		// data["id"] = value.Id
		// data["name"] = value.NickName
		// data["birthday"] = value.BirthDay
		// data["gender"] = value.Gender
		// data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(ctx *gin.Context) {
	// 表单验证
	var req forms.PassWordLoginForm
	if err := ctx.ShouldBind(&req); err != nil {
		// 返回错误信息
		HandlerValidatorError(ctx, err)
		return
	}

	if !store.Verify(req.CaptchaId, req.CaptchaAns, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 登录业务逻辑
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登陆失败",
				})
			}
			return
		}
	} else {
		// 以上只是查询到用户，未检验密码
		passRsp, err := global.UserSrvClient.CheckPassWord(ctx, &proto.PasswordCheckInfo{
			PassWord:          req.PassWord, // 请求参数输入的密码
			EncryptedPassWord: rsp.PassWord, // 调用grpc的服务返回的查询到的密码
		})
		if err != nil {
			zap.S().Errorw("[userSrvClient.CheckPassWord] 登录【验证密码服务错误】",
				"msg", err.Error(),
			)
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else {
			if passRsp.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := jwtClaims.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: &jwt.StandardClaims{
						NotBefore: time.Now().Unix(), // 签名的生效时间
						ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime,
						Issuer:    "lucien",
					},
				}
				accessToken, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}
				// 生成 refresh_token，有更长的有效期
				refreshClaims := claims
				refreshClaims.StandardClaims.ExpiresAt = time.Now().Unix() + 7*24*3600 // 7天
				refreshToken, err := j.CreateToken(refreshClaims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "生成refresh_token失败"})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"msg":           "登录成功",
					"id":            rsp.Id,
					"nick_name":     rsp.NickName,
					"access_token":  accessToken,
					"refresh_token": refreshToken,
					"expired_at":    (time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime) * 1000,
				})
			} else {
				zap.S().Errorw("[userSrvClient.CheckPassWord] 登录【密码验证错误】")
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"msg": "登录失败",
				})
			}
		}
	}
}

// 刷新token接口
func RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"  binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	j := middlewares.NewJWT()
	claims, err := j.ParseToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "refresh_token无效"})
		return
	}

	// 检查 refresh_token 是否过期
	if claims.ExpiresAt < time.Now().Unix() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "refresh_token已过期"})
		return
	}

	// 生成新的 access_token
	claims.StandardClaims.ExpiresAt = time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime
	accessToken, err := j.CreateToken(*claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "生成access_token失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expired_at":   (time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime) * 1000,
	})
}

// 用户注册
func Register(ctx *gin.Context) {
	var req forms.RegisterForm
	if err := ctx.ShouldBind(&req); err != nil {
		HandlerValidatorError(ctx, err)
		return
	}

	// 从redis中获取验证码
	value, err := global.Rdb.Get(context.Background(), req.Mobile).Result()
	if err == redis.Nil {
		zap.S().Errorw("验证码错误",
			"msg", err.Error(),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
	} else {
		if value != req.Code {
			zap.S().Errorw("验证码错误",
				"msg", err.Error(),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
		}
	}

	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: req.Mobile,
		PassWord: req.PassWord,
		Mobile:   req.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[CreateUser] 创建【注册用户失败:%s】", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// 创建用户成功，直接进入登录后的状态
	// 生成token
	j := middlewares.NewJWT()
	claims := jwtClaims.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: &jwt.StandardClaims{
			NotBefore: time.Now().Unix(), // 签名的生效时间
			ExpiresAt: time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime,
			Issuer:    "lucien",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":        "登录成功",
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime) * 1000,
	})
}

// 更新用户信息
func UpdateUser(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		HandlerValidatorError(ctx, err)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*jwtClaims.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "个人信息修改成功",
	})
}
