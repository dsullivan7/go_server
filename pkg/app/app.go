package app

import (
	"fmt"
	"go_server/internal/auth/auth0"
	goServerAlpaca "go_server/internal/broker/alpaca"
	"go_server/internal/config"
	"go_server/internal/db"
	goServerHTTP "go_server/internal/http"
	goServerZapLogger "go_server/internal/logger/zap"
	goServerPlaid "go_server/internal/plaid"
	"go_server/internal/server"
	"go_server/internal/services"
	"go_server/internal/store/gorm"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/plaid/plaid-go/plaid"

	"github.com/go-chi/chi"
)

const callerSkip = 2

func Run() {
	config, configErr := config.NewConfig()

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
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
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

	// initialize http client
	httpClient := goServerHTTP.NewClient()

	// initialize plaid
	plaidConfig := plaid.NewConfiguration()
	plaidConfig.AddDefaultHeader("PLAID-CLIENT-ID", config.PlaidClientID)
	plaidConfig.AddDefaultHeader("PLAID-SECRET", config.PlaidSecret)
	plaidConfig.UseEnvironment(plaid.Sandbox)
	plaidAPIClient := plaid.NewAPIClient(plaidConfig)
	plaidClient := goServerPlaid.NewClient(plaidAPIClient, config.PlaidRedirectURI)

	// initialize alpaca
	broker := goServerAlpaca.NewBroker(config.AlpacaAPIKey, config.AlpacaAPISecret, config.AlpacaAPIURL, httpClient)

	auth := auth0.NewAuth(config.Auth0Domain, config.Auth0Audience, logger)
	auth.Init()

	router := chi.NewRouter()
	handler := server.NewChiServer(config, router, srvc, store, plaidClient, broker, auth, logger)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: handler.Init(),
	}

	logger.Info(fmt.Sprintf("started on port: %s", config.Port))
	log.Fatal(httpServer.ListenAndServe())
}
