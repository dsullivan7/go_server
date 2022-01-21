package server

import (
	"net/http"

	"go_server/internal/auth"
	"go_server/internal/config"
	"go_server/internal/crawler"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/controllers"
	"go_server/internal/server/graph"
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
	resolver    *graph.Resolver
	middlewares *middlewares.Middlewares
	logger      logger.Logger
}

func NewChiServer(
	config *config.Config,
	router *chi.Mux,
	store store.Store,
	crawler crawler.Crawler,
	plaidClient plaid.Client,
	auth auth.Auth,
	logger logger.Logger,
) Server {
	utils := utils.NewServerUtils(logger)
	controllers := controllers.NewControllers(config, store, crawler, plaidClient, utils, logger)
	resolver := graph.NewResolver(config, store, logger)
	middlewares := middlewares.NewMiddlewares(config, store, auth, utils, logger)

	return &ChiServer{
		router:      router,
		config:      config,
		logger:      logger,
		controllers: controllers,
		resolver:    resolver,
		middlewares: middlewares,
	}
}
