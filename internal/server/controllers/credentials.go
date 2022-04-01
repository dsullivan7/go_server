package controllers

import (
	"encoding/json"
	"fmt"
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

	// encrypt username
	usernameEnc, errEnc := c.cipher.Encrypt(credentialPayload.Username, c.config.EncryptionKey)
	if errEnc != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEnc})

		return
	}

	credentialPayload.Username = usernameEnc

	// encrypt password
	passwordEnc, errEnc := c.cipher.Encrypt(credentialPayload.Password, c.config.EncryptionKey)
	if errEnc != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEnc})

		return
	}

	credentialPayload.Password = passwordEnc

	c.logger.Info(fmt.Sprint("value 1: ", usernameEnc))
	c.logger.Info(fmt.Sprint("value 2: ", passwordEnc))

	tag, err := c.store.CreateCredential(credentialPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, tag)
}
