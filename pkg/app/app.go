package app

import (
	"fmt"
	"go_server/internal/auth/auth0"
	"go_server/internal/config"
	"go_server/internal/db"
	goServerZapLogger "go_server/internal/logger/zap"
	"go_server/internal/server"
	"go_server/internal/store/gorm"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

func Run() {
	config, configErr := config.NewConfig()

	if configErr != nil {
		log.Fatal(configErr)
	}

	zapLogger, errZap := zap.NewProduction()

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

	auth := auth0.NewAuth(config.Auth0Domain, config.Auth0Audience)
	auth.Init()

	router := chi.NewRouter()

	handler := server.NewChiServer(config, router, store, auth, logger)

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: handler.Init(),
	}

	log.Fatal(httpServer.ListenAndServe())
}
