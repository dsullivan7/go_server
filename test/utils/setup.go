package utils

import (
	"go_server/internal/captcha/twocaptcha"
	"go_server/internal/config"
	goServerRodCrawler "go_server/internal/crawler/rod"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	goServerGormStore "go_server/internal/store/gorm"
	"go_server/test/mocks/auth"
	"go_server/test/mocks/plaid"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/go-rod/rod"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SetupUtility struct {
}

func NewSetupUtility() *SetupUtility {
	return &SetupUtility{}
}

func (setupUtility *SetupUtility) SetupIntegration() (*httptest.Server, *gorm.DB, DatabaseUtility, error) {
	config, configError := config.NewConfig()

	if configError != nil {
		return nil, nil, nil, configError
	}

	zapLogger, errZap := zap.NewProduction()

	if errZap != nil {
		return nil, nil, nil, errZap
	}

	logger := goServerZapLogger.NewLogger(zapLogger)

	connection, errConnection := db.NewSQLConnection(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)

	if errConnection != nil {
		return nil, nil, nil, errConnection
	}

	db, errDatabase := db.NewGormDB(connection)

	if errDatabase != nil {
		return nil, nil, nil, errDatabase
	}

	dbUtility := NewSQLDatabaseUtility(connection)

	store := goServerGormStore.NewStore(db)

	router := chi.NewRouter()

	authMock := auth.NewMockAuth()

	browser := rod.New()

	captchaKey := "key"

	captcha := twocaptcha.NewTwoCaptcha(captchaKey, logger)

	crawler := goServerRodCrawler.NewCrawler(browser, captcha)

	plaidClient := plaid.NewMockPlaidClient()

	handler := server.NewChiServer(config, router, store, crawler, plaidClient, authMock, logger)

	testServer := httptest.NewServer(handler.Init())

	return testServer, db, dbUtility, nil
}
