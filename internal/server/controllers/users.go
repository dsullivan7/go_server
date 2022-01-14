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
	userID := c.utils.GetURLParamUUID(r, "user_id")

	if (userID == uuid.Nil) {
		c.utils.HandleError(w, r, errors.HTTPNonExistentError{})

		return
	}

	user, err := c.store.GetUser(userID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, user)
}

func (c *Controllers) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	users, err := c.store.ListUsers(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, users)
}

func (c *Controllers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	errDecode := json.NewDecoder(r.Body).Decode(&userPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	user, err := c.store.CreateUser(userPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func (c *Controllers) ModifyUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	userID := c.utils.GetURLParamUUID(r, "user_id")

	errDecode := json.NewDecoder(r.Body).Decode(&userPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	user, err := c.store.ModifyUser(userID, userPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, user)
}

func (c *Controllers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "user_id")))

	err := c.store.DeleteUser(userID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
