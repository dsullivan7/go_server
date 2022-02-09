package utils

import (
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	"go_server/internal/server/graph"
	"go_server/internal/logger"
	"go_server/internal/auth"
	mockPlaid "go_server/test/mocks/plaid"
	mockBroker "go_server/test/mocks/broker"
	mockStore "go_server/test/mocks/store"
	mockAuth "go_server/test/mocks/auth"

	"github.com/go-chi/chi"
	"github.com/go-rod/rod"
	"go.uber.org/zap"
)

type TestServer struct {
	Server	server.Server
	Router	*chi.Mux
	Config      *config.Config
	Resolver    *graph.Resolver
	Logger      logger.Logger
	Store	*mockStore.MockStore
	Auth      	auth.Auth
	PlaidClient      	*mockPlaid.MockPlaid
	Broker      	*mockBroker.MockBroker
}

func NewTestServer() (*TestServer, error) {
	testServer := TestServer{}

	config, configError := config.NewConfig()

	if configError != nil {
		return nil, configError
	}

	testServer.Config = config

	zapLogger, errZap := zap.NewProduction()

	if errZap != nil {
		return nil, errZap
	}

	logger := goServerZapLogger.NewLogger(zapLogger)

	testServer.Logger = logger

	str := mockStore.NewMockStore()
	testServer.Store = str

	router := chi.NewRouter()
	testServer.Router = router

	ath := mockAuth.NewMockAuth()
	testServer.Auth = ath

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	pld := mockPlaid.NewMockPlaid()
	testServer.PlaidClient = pld

	brkr := mockBroker.NewMockBroker()
	testServer.Broker = brkr

	srvr := server.NewChiServer(config, router, str, crawler, pld, brkr, ath, logger)
	srvr.Init()

	testServer.Server = srvr

	return &testServer, nil
}
