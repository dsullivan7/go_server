package middlewares

import (
	"net/http"
  "context"

  "github.com/dgrijalva/jwt-go"

  "go_server/internal/services"
  "go_server/internal/models"
)

func User(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    auth0Id := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
    users := services.ListUsers(&models.User{ Auth0ID: auth0Id })
    if (len(users) > 0) {
      // Store the user making this request in the userModel field
      newContext := context.WithValue(r.Context(), "userModel", users[0])
      h.ServeHTTP(w, r.WithContext(newContext))
    } else {
      // No user found for this request
      h.ServeHTTP(w, r)
    }
  })
}