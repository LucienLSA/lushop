// pkg/captcha/captcha.go
package captcha

import (
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func Generate() (id, b64s, ans string, err error) {
	driver := base64Captcha.NewDriverDigit(
		80, 240, 4, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	return cp.Generate()
}

func Verify(id, answer string) bool {
	return store.Verify(id, answer, true)
}
