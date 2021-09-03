package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/errors"
	"net/http"

	jwtMiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

var errCert = fmt.Errorf("unable to find appropriate key")
var errAudience = fmt.Errorf("invalid audience")
var errIssuer = fmt.Errorf("invalid issuers")

func getPemCert(context context.Context, token *jwt.Token, domain string) (string, error) {
	cert := ""

	req, errRequest := http.NewRequestWithContext(
		context,
		http.MethodGet,
		fmt.Sprintf("https://%s/.well-known/jwks.json", domain),
		nil,
	)

	if errRequest != nil {
		return cert, fmt.Errorf("failed to create request: %w", errRequest)
	}

	res, errResponse := http.DefaultClient.Do(req)

	if errResponse != nil {
		return cert, fmt.Errorf("failed to request jwks endpoint: %w", errResponse)
	}
	defer res.Body.Close()

	var jwks = Jwks{}
	errDecode := json.NewDecoder(res.Body).Decode(&jwks)

	if errDecode != nil {
		return cert, fmt.Errorf("failed to decode jwks response: %w", errDecode)
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errCert
	}

	return cert, nil
}

func (m *Middlewares) Auth() func(http.Handler) http.Handler {
	jwtm := jwtMiddleware.New(jwtMiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := m.config.Auth0Audience

			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

			if !checkAud {
				return token, errAudience
			}
			// Verify 'iss' claim
			iss := fmt.Sprintf("https://%s/", m.config.Auth0Domain)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errIssuer
			}

			c := context.Background()
			cert, err := getPemCert(c, token, m.config.Auth0Domain)

			if err != nil {
				return nil, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
		// we will handle errors on return
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {},
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := jwtm.CheckJWT(w, r)
			if err != nil {
				m.utils.HandleError(w, r, errors.HTTPAuthError{Err: err})

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
