package crawler_test

import (
	"go_server/internal/captcha/twocaptcha"
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"

	"testing"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

func TestCrawler(t *testing.T) {
	t.Skip("No integration")

	browser := rod.New()

	zapLogger, _ := zap.NewProduction()

	logger := goServerZapLogger.NewLogger(zapLogger)

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	crawler.Login("https://www.connectebt.com/nyebtclient/siteLogonClient.recip", "username", "password")
}
