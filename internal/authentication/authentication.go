package authentication

import (
	"net/http"
)

type Authentication interface {
	CheckJWT(w http.ResponseWriter, r *http.Request) error
}
