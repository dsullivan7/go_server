package rod

import (
	"fmt"
	"go_server/internal/captcha"
	"go_server/internal/crawler"

	"time"

	"github.com/go-rod/rod"
)

const RenderWait = 5

type Crawler struct {
	browser *rod.Browser
	captcha captcha.Captcha
}

func NewCrawler(browser *rod.Browser, captcha captcha.Captcha) crawler.Crawler {
	return &Crawler{
		browser: browser,
		captcha: captcha,
	}
}

func (crawler *Crawler) Login(url string, username string, password string) string {
	crawler.browser.MustConnect()
	defer crawler.browser.MustClose()

	page := crawler.browser.MustPage(url)

	fr := page.MustElement("#main-iframe").MustFrame()

	googleKeyPointer := fr.MustElement(".g-recaptcha").MustAttribute("data-sitekey")
	googleKey := *googleKeyPointer

	captchaComplete, _ := crawler.captcha.SolveReCaptchaV2(googleKey, url)

	fr.MustEval(fmt.Sprintf("onCaptchaFinished('%s')", *captchaComplete))

	time.Sleep(RenderWait * time.Second)

	text := page.MustElement("body").MustText()

	return text
}
