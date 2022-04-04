package captcha

type Captcha interface {
	SolveReCaptchaV2(googleKey string, url string) (*string, error)
}
