package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go_server/internal/logger"

	jwtMiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

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

func getPemCert(token *jwt.Token, domain string) (string, error) {
	cert := ""

	req, errRequest := http.NewRequestWithContext(
		context.Background(),
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

type Auth struct {
	mw       *jwtMiddleware.JWTMiddleware
	domain   string
	audience string
	logger   logger.Logger
}

func NewAuth(domain string, audience string, logger logger.Logger) *Auth {
	logger.Info(fmt.Sprintf("Initialize auth with domain '%s' and audience '%s'", domain, audience))

	return &Auth{
		domain:   domain,
		audience: audience,
		logger:   logger,
	}
}

func (auth *Auth) Init() {
	mw := jwtMiddleware.New(jwtMiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := auth.audience

			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

			if !checkAud {
				return token, errAudience
			}
			// Verify 'iss' claim
			iss := fmt.Sprintf("https://%s/", auth.domain)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errIssuer
			}

			cert, err := getPemCert(token, auth.domain)

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

	auth.mw = mw
}

func (auth *Auth) CheckJWT(w http.ResponseWriter, r *http.Request) error {
	if err := auth.mw.CheckJWT(w, r); err != nil {
		return fmt.Errorf("error checking jwt: %w", err)
	}

	return nil
}
