package forms

type PassWordLoginForm struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式规范，自定义validator
	PassWord   string `form:"password" json:"password" binding:"required,min=3,max=10"`
	CaptchaAns string `form:"captcha_ans" json:"captcha_ans" binding:"required,min=5,max=5"`
	CaptchaId  string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式规范，自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=10"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}
