package rod

import (
	"fmt"
	"go_server/internal/crawler"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

type Crawler struct {
	browser *rod.Browser
}

func NewCrawler(browser *rod.Browser) crawler.Crawler {
	return &Crawler{
		browser: browser,
	}
}

func (crawler *Crawler) Login(url string, username string, password string) {
	crawler.browser.MustConnect()

	page := crawler.browser.MustPage(url)

	fr := page.MustElement("#main-iframe").MustFrame()

	googleKeyPointer := fr.MustElement(".g-recaptcha").MustAttribute("data-sitekey")
	googleKey := *googleKeyPointer

	captchaKey := "6da1d998757220665e090850725519bd"
	captchaUrl := fmt.Sprintf("http://2captcha.com/in.php?key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s", captchaKey, googleKey, url)

	resp1, _ := http.Get(captchaUrl)
	body1, _ := ioutil.ReadAll(resp1.Body)

	println(string(body1))

	captchaArray := strings.Split(string(body1), "|")
	captchaID := captchaArray[1]

	// captchaID := "68880267324"

	time.Sleep(120 * time.Second)
	captchaCompleteURL := fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=get&id=%s", captchaKey, captchaID)

	resp2, _ := http.Get(captchaCompleteURL)
	body2, _ := ioutil.ReadAll(resp2.Body)

	println(string(body2))

	captchaCompleteArray := strings.Split(string(body2), "|")
	captchaComplete := captchaCompleteArray[1]

	fr.MustEval(fmt.Sprintf("onCaptchaFinished('%s')", captchaComplete))
	// fr.MustEval("console.log('Hello!')")

	time.Sleep(5 * time.Second)
	text := page.MustElement("body").MustText()

	println("text")
	println(text)

	defer crawler.browser.MustClose()
}
