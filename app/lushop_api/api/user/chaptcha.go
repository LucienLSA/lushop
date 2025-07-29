package user

import (
	"lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// 获取图形验证码
// TODO 将该验证码存储在redis，实现冷却和过期实现，进行防刷
func GetCaptcha(ctx *gin.Context) {
	//表单验证
	captchaMobileForm := forms.CaptchaMobileForm{}
	if err := ctx.ShouldBind(&captchaMobileForm); err != nil {
		base.HandleValidatorError(ctx, err)
		return
	}
	if captchaMobileForm.Mobile != "" {
		coolKey := "captcha_cooldown:" + captchaMobileForm.Mobile
		exist, _ := global.RedisClient.Exists(ctx, coolKey).Result()
		if exist == 1 {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"msg": "请求太频繁，请稍后再试",
			})
			return
		}
		// 设置冷却时间60秒
		global.RedisClient.Set(ctx, coolKey, 1, 60*time.Second)
	}
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	//通过设置的driver放到自带的store
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, ans, err := cp.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误,:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 存储到 Redis，设置5分钟过期
	err = global.RedisClient.Set(ctx, "captcha:"+id, ans, 5*time.Minute).Err()
	if err != nil {
		zap.S().Errorf("验证码存入Redis失败: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码存储失败",
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
