package middlewares

import (
	"context"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"go_server/internal/server/consts"

	"github.com/go-chi/chi"
)

func (m *Middlewares) URLParam(param string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			paramValue := chi.URLParam(r, param)

			if paramValue == "me" {
				userModel := r.Context().Value(consts.UserModelKey)
				if userModel == nil {
					m.utils.HandleError(w, r, errors.HTTPServerError{Err: errAuth0SubString})

					return
				}

				// overwrite existing route context to set the path param
				routeContext := chi.RouteContext(r.Context())
				routeContext.URLParams.Add(param, userModel.(models.User).UserID.String())
				newContext := context.WithValue(r.Context(), chi.RouteCtxKey, routeContext)
				next.ServeHTTP(w, r.WithContext(newContext))

				return
			}

			// No param found
			next.ServeHTTP(w, r)
		})
	}
}
