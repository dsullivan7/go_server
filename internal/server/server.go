package server

import (
	"net/http"

	"go_server/internal/auth"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/server/controllers"
	"go_server/internal/server/middlewares"
	"go_server/internal/server/utils"
	"go_server/internal/store"

	"github.com/go-chi/chi"
)

type Server interface {
	Init() http.Handler
}

type ChiServer struct {
	router      *chi.Mux
	config      *config.Config
	controllers *controllers.Controllers
	middlewares *middlewares.Middlewares
	logger      logger.Logger
}

func NewChiServer(
	config *config.Config,
	router *chi.Mux,
	store store.Store,
	auth auth.Auth,
	logger logger.Logger,
) Server {
	utils := utils.NewServerUtils(logger)
	controllers := controllers.NewControllers(config, store, utils, logger)
	middlewares := middlewares.NewMiddlewares(config, store, auth, utils, logger)

	return &ChiServer{
		router:      router,
		config:      config,
		logger:      logger,
		controllers: controllers,
		middlewares: middlewares,
	}
}
