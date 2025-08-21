package forms

type CaptchaMobileForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
}

type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=10"`
	// CaptchaAns string `form:"captcha_ans" json:"captcha_ans" binding:"required,min=5,max=5"`
	// CaptchaId  string `form:"captcha_id" json:"captcha_id" binding:"required"`
}
type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式有规范可寻
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=10"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"`
}
type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=2,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}

type RefreshTokenForm struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

type UpdatePasswordForm struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=20"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}
