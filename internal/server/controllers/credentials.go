package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

func (c *Controllers) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var credentialPayload models.Credential

	errDecode := json.NewDecoder(r.Body).Decode(&credentialPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	// encrypt password
	credentialPayload.Password = c.cipher.Encrypt(credentialPayload.Password)

	tag, err := c.store.CreateCredential(credentialPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, tag)
}
