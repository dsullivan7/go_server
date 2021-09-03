package auth

import (
	"net/http"
)

type Auth interface {
	CheckJWT(w http.ResponseWriter, r *http.Request) error
}
