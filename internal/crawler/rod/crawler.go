package colly

import (
	"go_server/internal/crawler"

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

	text := page.MustElement("body").MustText()

	println("text")
	println(text)

	defer crawler.browser.MustClose()
}
