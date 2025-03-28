package api

import (
	"fmt"
	"lushopapi/user_web/forms"
	"lushopapi/user_web/global"

	"math/rand"
	"net/http"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
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
		HandlerValidatorError(ctx, err)
		return
	}
	config := &openapi.Config{
		// 您的AccessKey ID
		// AccessKeyId: tea.String(global.ServerConfig.AliSmsInfo.ApiKey),
		AccessKeyId: tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiKey)),
		// 您的AccessKey Secret
		// AccessKeySecret: tea.String(global.ServerConfig.AliSmsInfo.ApiSecrect),
		AccessKeySecret: tea.String(global.GetEnvInfoStr(global.ServerConfig.AliSmsInfo.ApiSecrect)),
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
	global.Rdb.Set(global.Rctx, req.Mobile, smsCode, time.Duration(global.ServerConfig.AliSmsInfo.Expire)*time.Second)
	// 后面注册时会将短信验证码带回来注册
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
