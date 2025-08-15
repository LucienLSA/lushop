package user

import (
	"context"
	"fmt"
	v2base "lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	"math/rand"
	"net/http"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 生成指定长度的验证码
func GenerateSmsCode(witdh int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 通过系统的环境变量获取
// 获取AccessKey ID
// 获取AccessKey Secret

func SendSmsAli(ctx *gin.Context) {
	var req forms.SendSmsForm
	if err := ctx.ShouldBind(&req); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}

	// 检查冷却时间
	cooldownKey := "cooldown:" + req.Mobile
	cooldownDuration := global.ServerConfig.AliSmsInfo.CoolDown

	if global.RedisClient.Exists(context.Background(), cooldownKey).Val() > 0 {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后再试",
		})
		return
	}

	// fmt.Println("ApiSecret:", global.ServerConfig.AliSmsInfo.ApiSecret)
	// fmt.Println("GetEnvInfoStr(ApiSecret):", global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiSecret))
	config := &openapi.Config{
		// 您的AccessKey ID
		// AccessKeyId: tea.String(global.ServerConfig.AliSmsInfo.ApiKey),
		AccessKeyId: tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiKey)),
		// 您的AccessKey Secret
		// AccessKeySecret: tea.String(global.ServerConfig.AliSmsInfo.ApiSecrect),
		AccessKeySecret: tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiSecret)),
		RegionId:        tea.String(global.ServerConfig.AliSmsInfo.RegionId),
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	client, _ := dysmsapi.NewClient(config)
	request := &dysmsapi.SendSmsRequest{}
	request.SetTemplateCode(global.ServerConfig.AliSmsInfo.TemplateCode)
	smsCode := GenerateSmsCode(6)
	request.SetTemplateParam("{\"code\":" + smsCode + "}")
	// 该参数值为假设值，请您根据实际情况进行填写
	request.SetPhoneNumbers(req.Mobile)
	// 该参数值为假设值，请您根据实际情况进行填写
	request.SetSignName(global.ServerConfig.AliSmsInfo.SignName)
	// fmt.Println(global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect,
	// 	global.ServerConfig.AliSmsInfo.TemplateCode, global.ServerConfig.AliSmsInfo.SignName)
	// fmt.Println(request.TemplateCode, request.TemplateParam,
	// 	request.PhoneNumbers, request.SignName)
	response, err := client.SendSms(request)
	fmt.Println(response)
	if err != nil {
		zap.S().Panic("调用发送阿里云短信服务失败", err.Error())
		return
	}

	// 将验证码保存起来，将手机号码作为redis的变量保存起来
	codeKey := "code:" + req.Mobile
	expireDuration := global.ServerConfig.AliSmsInfo.Expire
	// 设置过期时间
	global.RedisClient.Set(context.Background(), codeKey, smsCode, time.Duration(expireDuration)*time.Second)
	// 设置冷却时间
	global.RedisClient.Set(context.Background(), cooldownKey, "1", time.Duration(cooldownDuration)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

// func RefreshSmsAli(ctx *gin.Context) {
// 	var req forms.SendSmsForm
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		v2base.HandleValidatorError(ctx, err)
// 		return
// 	}
// 	// 检查冷却时间
// 	cooldownKey := "cooldown:" + req.Mobile
// 	cooldownDuration := global.ServerConfig.AliSmsInfo.CoolDown

// 	if global.RedisClient.Exists(context.Background(), cooldownKey).Val() > 0 {
// 		ctx.JSON(http.StatusTooManyRequests, gin.H{
// 			"msg": "请求过于频繁，请稍后再试",
// 		})
// 		return
// 	}

// 	config := &openapi.Config{
// 		AccessKeyId:     tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiKey)),
// 		AccessKeySecret: tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiSecret)),
// 		RegionId:        tea.String(global.ServerConfig.AliSmsInfo.RegionId),
// 	}
// 	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
// 	client, _ := dysmsapi.NewClient(config)

// 	request := &dysmsapi.SendSmsRequest{}
// 	request.SetTemplateCode(global.ServerConfig.AliSmsInfo.TemplateCode)
// 	smsCode := GenerateSmsCode(6)
// 	request.SetTemplateParam("{\"code\":" + smsCode + "}")
// 	request.SetPhoneNumbers(req.Mobile)
// 	request.SetSignName(global.ServerConfig.AliSmsInfo.SignName)

// 	response, err := client.SendSms(request)
// 	fmt.Println(response)
// 	if err != nil {
// 		zap.S().Panic("调用发送阿里云短信服务失败", err.Error())
// 		return
// 	}
// 	// 更新Redis中的验证码
// 	global.RedisClient.Set(context.Background(), req.Mobile, smsCode, time.Duration(global.ServerConfig.AliSmsInfo.Expire)*time.Second)

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"msg": "刷新成功",
// 	})
// }

func VerifySmsAli(ctx *gin.Context) {
	var req forms.VerifySmsForm
	if err := ctx.ShouldBind(&req); err != nil {
		v2base.HandleValidatorError(ctx, err)
		return
	}

	// 从Redis获取验证码
	codeKey := "code:" + req.Mobile
	expectedCode, err := global.RedisClient.Get(context.Background(), codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg": "验证码不存在或已过期",
			})
		} else {
			zap.S().Panic("Redis操作失败", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器错误",
			})
		}
		return
	}

	// 验证码匹配
	if expectedCode == req.Code {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "验证成功",
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "验证码错误",
		})
	}
}
