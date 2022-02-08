package server_test

import (
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	"go_server/test/mocks/auth"
	"go_server/test/mocks/plaid"
	"go_server/test/mocks/broker"
	"go_server/test/mocks/store"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServer(tParent *testing.T) {
	tParent.Parallel()

	config, configError := config.NewConfig()

	assert.Nil(tParent, configError)

	zapLogger, errZap := zap.NewProduction()

	assert.Nil(tParent, errZap)

	logger := goServerZapLogger.NewLogger(zapLogger)

	str := store.NewMockStore()

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	plaidClientMock := plaid.NewMockPlaid()

	brkr := broker.NewMockBroker()

	rtr := chi.NewRouter()

	ath := auth.NewMockAuth()

	handler := server.NewChiServer(config, rtr, str, crawler, plaidClientMock, brkr, ath, logger)

	handler.Init()

	tParent.Run("Test Init", func(t *testing.T) {
		t.Parallel()

		tctx := chi.NewRouteContext()
		assert.True(t, rtr.Match(tctx, "GET", "/api/users"))
	})
}
