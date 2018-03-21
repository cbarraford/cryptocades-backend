package users

type input struct {
	BTCAddr     string `json:"btc_address"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Referrer    string `json:"referrer"`
	CaptchaCode string `json:"captcha_code"`
}
