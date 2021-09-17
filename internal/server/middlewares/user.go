package middlewares

import (
	"context"
	"fmt"
	"go_server/internal/errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"go_server/internal/server/consts"
)

var errAuth0SubString = fmt.Errorf("auth0 sub must be a string")

func (m *Middlewares) User() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth0Id, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)

			if !ok {
				m.utils.HandleError(w, r, errors.HTTPServerError{Err: errAuth0SubString})

				return
			}
			users, err := m.store.ListUsers(map[string]interface{}{"auth0_id": auth0Id})

			if err != nil {
				m.utils.HandleError(w, r, errors.HTTPServerError{Err: err})

				return
			}

			if len(users) == 1 {
				// Store the user making this request in the userModel field
				newContext := context.WithValue(r.Context(), consts.UserModelKey, users[0])
				next.ServeHTTP(w, r.WithContext(newContext))

				return
			}

			// No user found for this request
			next.ServeHTTP(w, r)
		})
	}
}
