package middlewares

import (
	"fmt"
	"encoding/json"
	"net/http"
	"errors"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N string `json:"n"`
	E string `json:"e"`
	X5c []string `json:"x5c"`
}

func getPemCert(token *jwt.Token, domain string) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", domain))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func (m *Middlewares) Auth() func(next http.Handler) http.Handler {
	return jwtmiddleware.New(jwtmiddleware.Options {
    ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
      // Verify 'aud' claim
      aud := m.config.Auth0Audience

      checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)

      if !checkAud {
        return token, errors.New("Invalid audience.")
      }
      // Verify 'iss' claim
      iss := fmt.Sprintf("https://%s/", m.config.Auth0Domain)
      checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
      if !checkIss {
        return token, errors.New("Invalid issuer.")
      }

      cert, err := getPemCert(token, m.config.Auth0Domain)

      if err != nil {
        return nil, err
      }

      result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
      return result, nil
    },
    SigningMethod: jwt.SigningMethodRS256,
  }).Handler
}