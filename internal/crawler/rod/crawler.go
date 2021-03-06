package rod

import (
	"fmt"
	"go_server/internal/captcha"
	"go_server/internal/crawler"

	"context"
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
	cxt := crawler.LoginPage(url, username, password)

	return crawler.GetAuthKey(cxt, url)
}

func (crawler *Crawler) LoginPage(url string, username string, password string) context.Context {
	browser := rod.New()
	browser.MustConnect()

	defer browser.MustClose()

	page := browser.MustPage(url)

	page.MustWaitLoad()

	page.MustElement("#username").MustInput("dbsullivan23+test@gmail.com")
	page.MustElement("#password").MustInput("TestPassword")
	page.MustElement("[type=\"submit\"]").MustClick()

	page.MustWaitLoad()
	time.Sleep(RenderWait * time.Second)

	context := browser.GetContext()

	return context
}

func (crawler *Crawler) GetAuthKey(ctx context.Context, url string) string {
	browser := rod.New().Context(ctx)
	browser.MustConnect()

	defer browser.MustClose()

	time.Sleep(RenderWait * time.Second)

	page := browser.MustPage(url)

	page.MustWaitLoad()
	time.Sleep(RenderWait * time.Second)

	page.MustScreenshot("my.png")

	authKey := page.MustElement("textarea").MustText()

	return authKey
}

func (crawler *Crawler) LoginEBT(url string, username string, password string) string {
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
