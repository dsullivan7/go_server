package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"go_server/internal/auth"
	"go_server/test/mocks/consts"
)

type MockAuth struct{}

func NewMockAuth() auth.Auth {
	return &MockAuth{}
}

func (auth *MockAuth) CheckJWT(w http.ResponseWriter, r *http.Request) error {
	jwtToken := &jwt.Token{Claims: jwt.MapClaims{"sub": consts.LoggedInAuth0Id}}

	newRequest := r.WithContext(context.WithValue(r.Context(), "user", jwtToken)) //nolint:revive,staticcheck
	*r = *newRequest

	return nil
}
