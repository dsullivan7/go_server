package controllers

import (
	"go_server/internal/config"
	goServerErrors "go_server/internal/errors"
	"go_server/internal/logger"
	"go_server/internal/store"
	"net/http"

	"github.com/go-chi/render"
)

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


func (c *Controllers) HandlePanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var foundError error
				switch x := err.(type) {
	        case string:
            foundError = goServerErrors.RunTimeError{ ErrorText: x }
	        case error:
            foundError = x
	        default:
						foundError = goServerErrors.RunTimeError{ ErrorText: "unknown" }
	      }

				c.handleError(w, r, goServerErrors.HTTPServerError{Err: foundError})
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (c *Controllers) handleError(w http.ResponseWriter, r *http.Request, err goServerErrors.HTTPError) {
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
