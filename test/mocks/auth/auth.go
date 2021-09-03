package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (auth *Auth) CheckJWT(w http.ResponseWriter, r *http.Request) error {
	jwtToken := &jwt.Token{Claims: jwt.MapClaims{"sub": "auth0|loggedInUser"}}

	newRequest := r.WithContext(context.WithValue(r.Context(), "user", jwtToken)) //nolint:revive,staticcheck
	*r = *newRequest

	return nil
}
