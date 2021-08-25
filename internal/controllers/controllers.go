package controllers

import (
	"go_server/internal/config"
	"go_server/internal/errors"
	"go_server/internal/logger"
	"go_server/internal/store"
	"net/http"

	"github.com/go-chi/render"
)

const HTTP400 = 400

type Controllers struct {
	store  store.Store
	config *config.Config
	logger logger.Logger
}

func NewControllers(store store.Store, config *config.Config, logger logger.Logger) *Controllers {
	return &Controllers{
		store:  store,
		config: config,
		logger: logger,
	}
}

func (c *Controllers) handleError(w http.ResponseWriter, r *http.Request, err errors.HTTPError) {
	c.logger.ErrorWithMeta(
		"Error",
		map[string]interface{}{
			"err": err.GetError(),
		},
	)

	w.WriteHeader(err.GetHTTPStatus())

	logJSON := map[string]interface{}{
		"message": err.GetMessage(),
	}
	render.JSON(w, r, logJSON)
}
