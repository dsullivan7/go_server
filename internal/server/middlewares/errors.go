package middlewares

import (
	"go_server/internal/errors"
	"net/http"
)

func (m *Middlewares) HandlePanic() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					var foundError error
					switch x := err.(type) {
					case string:
						foundError = errors.RunTimeError{ErrorText: x}
					case error:
						foundError = x
					default:
						foundError = errors.RunTimeError{ErrorText: "unknown"}
					}

					m.utils.HandleError(w, r, errors.HTTPServerError{Err: foundError})
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
