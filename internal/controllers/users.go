package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (c *Controllers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user, err := c.store.GetUser(userID)

	if err != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, user)
}

func (c *Controllers) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	users, err := c.store.ListUsers(query)

	if err != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, users)
}

func (c *Controllers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	errDecode := json.NewDecoder(r.Body).Decode(&userPayload)
	if errDecode != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	// if (userPayload.Auth0ID == "") {
	// 	auth0Id := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	// 	userPayload.Auth0ID = auth0Id
	// }

	user, err := c.store.CreateUser(userPayload)

	if err != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func (c *Controllers) ModifyUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	errDecode := json.NewDecoder(r.Body).Decode(&userPayload)
	if errDecode != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	user, err := c.store.ModifyUser(userID, userPayload)

	if err != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, user)
}

func (c *Controllers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	err := c.store.DeleteUser(userID)

	if err != nil {
		c.handleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
