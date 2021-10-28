package crawler_test

import (
	goServerRodCrawler "go_server/internal/crawler/rod"

  "github.com/go-rod/rod"
  "testing"
)

func TestCrawler(t *testing.T) {
  browser := rod.New()

	crawler := goServerRodCrawler.NewCrawler(browser)

  crawler.Login("https://www.connectebt.com/nyebtclient/siteLogonClient.recip", "username", "password")
}
