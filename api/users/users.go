package users

type input struct {
	BTCAddr      string `json:"btc_address"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	ReferralCode string `json:"referral_code"`
	CaptchaCode  string `json:"captcha_code"`
}
