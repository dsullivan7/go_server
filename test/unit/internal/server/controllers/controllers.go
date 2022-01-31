package controllers

import (
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server/controllers"
	"go_server/internal/server/utils"
	"go_server/test/mocks/bank"
	"go_server/test/mocks/store"

	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

func Setup() (*controllers.Controllers, *store.MockStore, error) {
	config, errConfig := config.NewConfig()

	if errConfig != nil {
		return nil, nil, errConfig
	}

	zapLogger, errZap := zap.NewProduction()
	if errZap != nil {
		return nil, nil, errZap
	}

	logger := goServerZapLogger.NewLogger(zapLogger)

	store := store.NewMockStore()

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	bnk := bank.NewMockBank()

	utils := utils.NewServerUtils(logger)
	controllers := controllers.NewControllers(config, store, crawler, bnk, utils, logger)

	return controllers, store, nil
}
