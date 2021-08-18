package app

import (
	"go_server/internal/server"
	"go_server/internal/db"
	"go_server/internal/controllers"
	"go_server/internal/store"
	"go_server/internal/config"
	"go_server/internal/logger"
	"github.com/go-chi/chi"
)

func Run() {
	config := config.NewConfig()
	logger := logger.NewZapLogger()

  db, err := db.NewDatabase(
		config.DBHost,
		config.DBName,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBSSL,
	)

	store := store.NewGormStore(db)

	controllers := controllers.NewControllers(store, config, logger)

	router := chi.NewRouter()

	server := server.NewServer(router, controllers, config, logger)

	server.Run()
}
