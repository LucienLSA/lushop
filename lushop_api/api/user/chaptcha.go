package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// 获取图形验证码
func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	//通过设置的driver放到store里面
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, ans, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误,:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captcha_id":   id,
		"picture_path": b64s,
		"captcha_ans":  ans, // 这里做测试才返回
	})
	//store.Verify(id,b64s,true)
}
