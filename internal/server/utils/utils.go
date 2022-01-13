package utils

import (
	"go_server/internal/errors"
	"go_server/internal/logger"
	"go_server/internal/models"
	"go_server/internal/server/consts"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	}
	render.JSON(w, r, logJSON)
}

func (s *ServerUtils) GetURLParamUUID(r *http.Request, param string) uuid.UUID {
	paramValue := chi.URLParam(r, param)

	if paramValue == "me" {
		return r.Context().Value(consts.UserModelKey).(models.User).UserID
	}

	return uuid.Must(uuid.Parse(paramValue))
}

func (s *ServerUtils) GetQueryParamUUID(r *http.Request, param string) uuid.UUID {
	paramValue := r.URL.Query().Get(param)

	// check if the query parameter is not specified
	if paramValue == "" {
		return uuid.Nil
	}

	// check if the query parameter refers to the requesting user
	if paramValue == "me" {
		return r.Context().Value(consts.UserModelKey).(models.User).UserID
	}

	// return the parameter as specified
	return uuid.Must(uuid.Parse(paramValue))
}
