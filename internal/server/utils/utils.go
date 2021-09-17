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

func GetURLParamUUID(r *http.Request, param string) uuid.UUID {
	paramValue := chi.URLParam(r, param)

	var paramUUID uuid.UUID

	if paramValue == "me" {
		paramUUID = r.Context().Value(consts.UserModelKey).(models.User).UserID
	} else {
		paramUUID = uuid.Must(uuid.Parse(paramValue))
	}

	return paramUUID
}
