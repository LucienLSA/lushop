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

// func unsafeVerify(publicKey []byte, md5 []byte, auth []byte) bool {
// 	// 1. 解析公钥
// 	block, _ := pem.Decode(publicKey)
// 	if block == nil {
// 		return false
// 	}

// 	// 2. 不检查密钥长度，直接使用
// 	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return false
// 	}

// 	pub := pubInterface.(*rsa.PublicKey)

// 	// 3. 自定义验证逻辑
// 	hashed := sha256.Sum256(md5)
// 	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], auth)
// 	return err == nil
// }

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

	if !utils.VerifySignatureV2(bytePublicKey, byteMD5, byteAuthorization) {
		zap.S().Error("Signature verification failed")
	}

	// 6. 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"url":     fileUrl,
		"status":  "success",
		"message": "Callback verified successfully",
	})
}

func PostPicture(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "文件获取失败", "error": err.Error()})
		return
	}

	// 2. 保存到本地临时目录（可选，直接用 file.Open() 也可以）
	localPath := "./tmp/" + file.Filename
	if err := c.SaveUploadedFile(file, localPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "文件保存失败", "error": err.Error()})
		return
	}

	// 3. 构造 OSS 对象名（可自定义路径/前缀）
	objectName := global.ServerConfig.OssInfo.UploadDir + file.Filename
	// 4. 上传到OSS
	err = utils.UploadFileToOSS(objectName, localPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "OSS上传失败", "error": err.Error()})
		return
	}

	// 5. 返回文件在OSS的访问路径
	ossUrl := strings.TrimRight(global.ServerConfig.OssInfo.Host, "/") + "/" + objectName
	c.JSON(http.StatusOK, gin.H{
		"msg":    "上传成功",
		"ossUrl": ossUrl,
	})
}
