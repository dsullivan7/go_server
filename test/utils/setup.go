package utils

import (
	"go_server/internal/authentication"
	"go_server/internal/config"
	"go_server/internal/gov"
	"go_server/internal/logger"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	"go_server/internal/server/graph"
	mockAuthentication "go_server/test/mocks/authentication"
	mockBank "go_server/test/mocks/bank"
	mockBroker "go_server/test/mocks/broker"
	mockCipher "go_server/test/mocks/cipher"
	mockPlaid "go_server/test/mocks/plaid"
	mockServices "go_server/test/mocks/services"
	mockStore "go_server/test/mocks/store"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type TestServer struct {
	Server         server.Server
	Router         *chi.Mux
	Config         *config.Config
	Service        *mockServices.MockService
	Resolver       *graph.Resolver
	Logger         logger.Logger
	Store          *mockStore.MockStore
	Authentication authentication.Authentication
	PlaidClient    *mockPlaid.MockPlaid
	Bank           *mockBank.MockBank
	Broker         *mockBroker.MockBroker
	Cipher         *mockCipher.MockCipher
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

	ath := mockAuthentication.NewMockAuthentication()
	testServer.Authentication = ath

	pld := mockPlaid.NewMockPlaid()
	testServer.PlaidClient = pld

	srvc := mockServices.NewMockService()
	testServer.Service = srvc

	brkr := mockBroker.NewMockBroker()
	testServer.Broker = brkr

	bnk := mockBank.NewMockBank()
	testServer.Bank = bnk

	gv := gov.NewGov()

	cphr := mockCipher.NewMockCipher()
	testServer.Cipher = cphr

	srvr := server.NewChiServer(config, router, srvc, str, pld, brkr, bnk, gv, cphr, ath, logger)
	srvr.Init()

	testServer.Server = srvr

	return &testServer, nil
}
