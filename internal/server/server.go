package server

import (
	"net/http"

	"go_server/internal/authentication"
	"go_server/internal/broker"
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/plaid"
	"go_server/internal/server/controllers"
	"go_server/internal/server/graph"
	"go_server/internal/server/middlewares"
	"go_server/internal/server/utils"
	"go_server/internal/services"
	"go_server/internal/store"

	"github.com/go-chi/chi"
)

type Server interface {
	Init() http.Handler
	GetControllers() *controllers.Controllers
	GetMiddlewares() *middlewares.Middlewares
}

type ChiServer struct {
	controllers *controllers.Controllers
	middlewares *middlewares.Middlewares
	router      *chi.Mux
	config      *config.Config
	resolver    *graph.Resolver
	logger      logger.Logger
}

func NewChiServer(
	cfg *config.Config,
	router *chi.Mux,
	srvc services.IService,
	str store.Store,
	pld plaid.IClient,
	brkr broker.Broker,
	ath authentication.Authentication,
	lggr logger.Logger,
) Server {
	utils := utils.NewServerUtils(lggr)
	controllers := controllers.NewControllers(cfg, str, srvc, pld, brkr, utils, lggr)
	resolver := graph.NewResolver(cfg, str, lggr)
	middlewares := middlewares.NewMiddlewares(cfg, str, ath, utils, lggr)

	return &ChiServer{
		controllers: controllers,
		middlewares: middlewares,
		router:      router,
		config:      cfg,
		logger:      lggr,
		resolver:    resolver,
	}
}

func (s *ChiServer) GetControllers() *controllers.Controllers {
	return s.controllers
}

func (s *ChiServer) GetMiddlewares() *middlewares.Middlewares {
	return s.middlewares
}
