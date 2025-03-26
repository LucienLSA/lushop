package api

import (
	"context"
	"fmt"
	"lushopapi/user_web/forms"
	"lushopapi/user_web/global"
	"lushopapi/user_web/global/response"
	"lushopapi/user_web/global/types"
	"lushopapi/user_web/middlewares"
	"lushopapi/user_web/proto"
	"lushopapi/user_web/utils/jwtClaims"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
					"msg": "用户服务不可用",
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
	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
	// 拨号连接grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "5")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
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

	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
	// 拨号连接grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
			"msg", err.Error(),
		)
	}
	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 登录业务逻辑
	rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
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
		// 只是查询到用户，未检验密码
		passRsp, err := userSrvClient.CheckPassWord(ctx, &proto.PasswordCheckInfo{
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
					StandardClaims: jwt.StandardClaims{
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
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + global.ServerConfig.JwtInfo.ExpireTime) * 1000,
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
