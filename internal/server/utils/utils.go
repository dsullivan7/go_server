package utils

import (
	"go_server/internal/errors"
	"go_server/internal/logger"
	"net/http"

	"github.com/go-chi/render"
)

type ServerUtils struct {
	logger logger.Logger
}

func NewServerUtils(
	logger logger.Logger,
) *ServerUtils {
	return &ServerUtils{
		logger: logger,
	}
}

func (s *ServerUtils) HandleError(w http.ResponseWriter, r *http.Request, err errors.HTTPError) {
	s.logger.ErrorWithMeta(
		"Error",
		map[string]interface{}{
			"err": err.GetError(),
		},
	)

	w.WriteHeader(err.GetHTTPStatus())

	logJSON := map[string]interface{}{
		"message": err.GetMessage(),
		"code":    err.GetCode(),
	}
	render.JSON(w, r, logJSON)
}
