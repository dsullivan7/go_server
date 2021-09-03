package middlewares

import (
	"go_server/internal/errors"
	"net/http"
)

func (m *Middlewares) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := m.auth.CheckJWT(w, r)
			if err != nil {
				m.utils.HandleError(w, r, errors.HTTPAuthError{Err: err})

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
