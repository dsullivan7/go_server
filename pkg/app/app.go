package app

import (
	"go_server/internal/config"
	"go_server/internal/controllers"
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

	driver, errDriver := db.NewDatabaseDiver(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)

	if errDriver != nil {
		log.Fatal(errDriver)
	}

	db, errDatabase := db.NewDatabase(driver)

	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	store := store.NewGormStore(db)

	controllers := controllers.NewControllers(store, config, logger)

	router := chi.NewRouter()

	server := server.NewServer(router, controllers, config, logger)

	server.Run()
}
