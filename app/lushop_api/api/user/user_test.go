package user

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// // 测试设置 - 在每个测试前初始化gin和必要的依赖
// func setupTest(t *testing.T) *gin.Engine {

// }

// // TestRegister 测试用户注册接口
// func TestRegister(t *testing.T) {

// }

// // TestLogout 测试登出接口
// func TestLogout(t *testing.T) {

// }

// // TestUpdatePassword 测试修改密码接口
// func TestUpdatePassword(t *testing.T) {

// }

// // TestGetUserDetail 测试获取用户详情接口
// func TestGetUserDetail(t *testing.T) {

// }

// // TestGetUserList 测试管理员用户列表接口
// func TestGetUserList(t *testing.T) {

// }

// // TestUpdateUserInfo 测试更新用户信息接口
// func TestUpdateUserInfo(t *testing.T) {

// }

// // TestGetCaptchaV2 测试获取图形验证码接口
// func TestGetCaptchaV2(t *testing.T) {

// }

// // TestVerifyCaptcha 测试验证图形验证码接口
// func TestVerifyCaptcha(t *testing.T) {

// }

// // TestSendSmsAli 测试发送手机验证码接口
// func TestSendSmsAli(t *testing.T) {
// 	router := setupTest(t)
// 	router.POST("/captcha/sms/send", SendSmsAli)

// 	inputBody := `{
// 		"mobile": "13800138000"
// 	}`

// 	req, _ := http.NewRequest("POST", "/captcha/sms/send", bytes.NewBufferString(inputBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }

// // TestVerifySmsAli 测试验证手机验证码接口
// func TestVerifySmsAli(t *testing.T) {
// 	router := setupTest(t)
// 	router.GET("/captcha/sms/verify", VerifySmsAli)

// 	req, _ := http.NewRequest("GET", "/captcha/sms/verify?mobile=13800138000&code=123456", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// }
