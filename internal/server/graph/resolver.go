package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"go_server/internal/config"
	"go_server/internal/logger"
	"go_server/internal/store"
)

type Resolver struct {
	config *config.Config
	store  store.Store
	logger logger.Logger
}

func NewResolver(
	config *config.Config,
	store store.Store,
	logger logger.Logger,
) *Resolver {
	return &Resolver{
		store:  store,
		config: config,
		logger: logger,
	}
}
