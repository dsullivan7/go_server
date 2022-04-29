package app

import (
	"fmt"
	"go_server/internal/authentication/auth0"
	goServerAlpaca "go_server/internal/broker/alpaca"
	"go_server/internal/cipher"
	"go_server/internal/config"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/plaid"
	"go_server/internal/bank/dwolla"
	"go_server/internal/server"
	"go_server/internal/services"
	"go_server/internal/gov"
	"go_server/internal/store/gorm"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

const callerSkip = 2

func Run() {
	cfg, configErr := config.NewConfig()

	if configErr != nil {
		log.Fatal(configErr)
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapLogger, errZap := zapConfig.Build(zap.AddCallerSkip(callerSkip))

	if errZap != nil {
		log.Fatal(errZap)
	}

	logger := goServerZapLogger.NewLogger(zapLogger)

	connection, errConnection := db.NewSQLConnection(
		cfg.DBHost,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBSSL,
	)

	if errConnection != nil {
		log.Fatal(errConnection)
	}

	db, errDatabase := db.NewGormDB(connection)

	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	store := gorm.NewStore(db)

	srvc := services.NewService()

	// initialize plaid
	plaidClient := plaid.NewClient(
		cfg.PlaidClientID,
		cfg.PlaidSecret,
		cfg.PlaidAPIURL,
		cfg.PlaidRedirectURI,
		logger,
	)

	// initialize alpaca
	broker := goServerAlpaca.NewBroker(cfg.AlpacaAPIKey, cfg.AlpacaAPISecret, cfg.AlpacaAPIURL)

	dwollaBank := dwolla.NewBank(
		cfg.DwollaAPIKey,
		cfg.DwollaAPISecret,
		cfg.DwollaAPIURL,
		cfg.DwollaWebhookURL,
		cfg.DwollaWebhookSecret,
    logger,
	)

	cphr := cipher.NewCipher()

	gv := gov.NewGov()

	auth := auth0.NewAuth(cfg.Auth0Domain, cfg.Auth0Audience, logger)
	auth.Init()

	router := chi.NewRouter()
	handler := server.NewChiServer(cfg, router, srvc, store, plaidClient, broker, dwollaBank, gv, cphr, auth, logger)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: handler.Init(),
	}

	logger.Info(fmt.Sprintf("started on port: %s", cfg.Port))
	log.Fatal(httpServer.ListenAndServe())
}
