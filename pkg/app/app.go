package app

import (
	"go_server/internal/config"
	"go_server/internal/db"
	"go_server/internal/logger"
	"go_server/internal/server"
	"go_server/internal/store"
	"log"

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

	logger := logger.NewZapLogger(zapLogger)

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

	store := store.NewGormStore(db)

	router := chi.NewRouter()

	server := server.NewServer(config, router, store, logger)

	server.Init()
	server.Run()
}
