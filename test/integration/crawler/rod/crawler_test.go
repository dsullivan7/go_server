package crawler_test

import (
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/captcha/twocaptcha"

	"testing"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

func TestCrawler(t *testing.T) {
	browser := rod.New()

	zapLogger, _ := zap.NewProduction()

	logger := goServerZapLogger.NewLogger(zapLogger)

	captchaKey := "6da1d998757220665e090850725519bd"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	crawler.Login("https://www.connectebt.com/nyebtclient/siteLogonClient.recip", "username", "password")
}
