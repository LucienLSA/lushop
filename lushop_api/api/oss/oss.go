package oss

import (
	"fmt"
	"lushopapi/global"
	"net/http"
	"net/url"
	"strings"

	"lushopapi/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Token(c *gin.Context) {
	response := utils.Get_policy_token()
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Origin", "*")
	c.String(200, response)
}

// 验证阿里云 OSS 在文件上传完成后发送的回调请求的合法性，并返回文件访问 URL
func HandlerRequest(ctx *gin.Context) {
	fmt.Println("\nHandle Post Request ... ")
	// Get PublicKey bytes
	bytePublicKey, err := utils.GetPublicKey(ctx)
	if err != nil {
		zap.S().Errorf("GetPublicKey err:%s", err.Error())
		utils.ResponseFailed(ctx)
		return
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := utils.GetAuthorization(ctx)
	if err != nil {
		zap.S().Errorf("GetAuthorization err:%s", err.Error())
		utils.ResponseFailed(ctx)
		return
	}

	// Get MD5 bytes from Newly Constructed Authrization String. 构建并计算 MD5 校验值
	byteMD5, bodyStr, err := utils.GetMD5FromNewAuthString(ctx)
	if err != nil {
		zap.S().Errorf("GetMD5FromNewAuthString err:%s", err.Error())
		utils.ResponseFailed(ctx)
		return
	}

	// 解析回调参数
	decodeurl, err := url.QueryUnescape(bodyStr)
	if err != nil {
		zap.S().Errorf("QueryUnescape err:%s", err.Error())
		fmt.Println(err)
	}
	fmt.Println(decodeurl)
	params := make(map[string]string)
	datas := strings.Split(decodeurl, "&")
	for _, v := range datas {
		sdatas := strings.Split(v, "=")
		fmt.Println(v)
		params[sdatas[0]] = sdatas[1]
	}
	fileName := params["filename"]
	zap.S().Infof("%s/%s", global.ServerConfig.OssInfo.Host, fileName)
	// 生成文件访问 URL
	fileUrl := fmt.Sprintf("%s/%s", global.ServerConfig.OssInfo.Host, fileName)

	// 签名验证与响应
	// verifySignature and response to client
	if utils.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want accoding to callback_body ...
		ctx.JSON(http.StatusOK, gin.H{
			"url": fileUrl,
		})
		utils.ResponseSuccess(ctx) // response OK : 200
	} else {
		utils.ResponseFailed(ctx) // response FAILED : 400
	}
}
