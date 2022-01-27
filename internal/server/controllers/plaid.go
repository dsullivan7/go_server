package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"net/http"

	"github.com/go-chi/render"
)

func (c *Controllers) CreatePlaidToken(w http.ResponseWriter, r *http.Request) {
	var tokenPayload map[string]string

	errDecode := json.NewDecoder(r.Body).Decode(&tokenPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	userID := tokenPayload["user_id"]

	plaidToken, err := c.plaidClient.CreatePlaidToken(userID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, render.M{"plaid_token": plaidToken})
}