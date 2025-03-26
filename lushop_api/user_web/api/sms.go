package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func SendSms(ctx *gin.Context) {
	// 后面注册时会将短信验证码带回来注册
	ctx.JSON(http.StatusOK, gin.H{})
	// 将验证码保存起来，将手机号码作为redis的变量保存起来

}
