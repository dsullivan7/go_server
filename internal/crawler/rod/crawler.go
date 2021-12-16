package rod

import (
	"fmt"
	"go_server/internal/crawler"
	"go_server/internal/captcha"
	"time"

	"github.com/go-rod/rod"
)

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

func (crawler *Crawler) Login(url string, username string, password string) {
	crawler.browser.MustConnect()

	page := crawler.browser.MustPage(url)

	fr := page.MustElement("#main-iframe").MustFrame()

	googleKeyPointer := fr.MustElement(".g-recaptcha").MustAttribute("data-sitekey")
	googleKey := *googleKeyPointer

	captchaComplete, _ := crawler.captcha.SolveReCaptchaV2(googleKey, url)

	fr.MustEval(fmt.Sprintf("onCaptchaFinished('%s')", *captchaComplete))

	time.Sleep(5 * time.Second)

	text := page.MustElement("body").MustText()

	println("text")
	println(text)

	defer crawler.browser.MustClose()
}
