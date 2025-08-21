package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}

type VerifySmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻
	Code   string `form:"code" json:"code" binding:"required,len=6"`      // 验证码，假设是6位
}
