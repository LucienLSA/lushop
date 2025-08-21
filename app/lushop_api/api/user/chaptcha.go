package user

import (
	"lushopapi/api/base"
	"lushopapi/forms"
	"lushopapi/global"
	"lushopapi/utils/captcha"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// 获取图形验证码

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
		// 不设置冷却时间
		global.RedisClient.Set(ctx, coolKey, 1, 1*time.Millisecond)
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

func GetCaptchaV2(ctx *gin.Context) {
	id, b64s, ans, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码错误,:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 存储到Redis（5分钟过期）这里存入captcha_id作为key，ans作为value
	if err := global.RedisClient.SetNX(ctx, "captcha:"+id, ans, 5*time.Minute).Err(); err != nil {
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
}

// // 刷新验证码
// func RefreshCaptcha(ctx *gin.Context) {
// 	captchaId := ctx.Query("captcha_id")
// 	if captchaId == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"msg": "验证码ID不能为空",
// 		})
// 		return
// 	}

// 	// 生成新的图形验证码
// 	id, b64s, ans, err := captcha.Generate()
// 	if err != nil {
// 		zap.S().Errorf("生成验证码错误,:", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"msg": "验证码错误",
// 		})
// 		return
// 	}

// 	// 删除旧的验证码
// 	global.RedisClient.Del(ctx, "captcha:"+captchaId)

// 	// 存储新的验证码到Redis，设置5分钟过期
// 	err = global.RedisClient.Set(ctx, "captcha:"+id, ans, 5*time.Minute).Err()
// 	if err != nil {
// 		zap.S().Errorf("验证码存入Redis失败: %s", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"msg": "验证码存储失败",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"captcha_id":   id,
// 		"picture_path": b64s,
// 		"captcha_ans":  ans,
// 	})
// }

// 验证验证码
func VerifyCaptcha(ctx *gin.Context) {
	captchaId := ctx.Query("captcha_id")
	captchaAns := ctx.Query("captcha_ans")
	if captchaId == "" || captchaAns == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码ID和答案不能为空",
		})
		return
	}

	// 从Redis获取验证码答案
	ans, err := global.RedisClient.Get(ctx, "captcha:"+captchaId).Result()
	if err != nil {
		zap.S().Errorf("获取验证码失败: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码已过期或不存在",
		})
		return
	}

	if ans != captchaAns {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	if !captcha.Verify(captchaId, captchaAns) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 验证成功后删除验证码 阅后即焚
	global.RedisClient.Del(ctx, "captcha:"+captchaId)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "验证成功",
	})
}
