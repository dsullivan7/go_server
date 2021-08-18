package controllers

import (
	"net/http"
	"encoding/json"

	"go_server/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"

  // "github.com/dgrijalva/jwt-go"
)

func (c *Controllers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user := c.store.GetUser(userID)

	render.JSON(w, r, user)
}

func (c *Controllers) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	users := c.store.ListUsers(query)

	render.JSON(w, r, users)
}

func (c *Controllers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	err := json.NewDecoder(r.Body).Decode(&userPayload)

	if (err != nil) {
		w.WriteHeader(400)
		return
	}

	// if (userPayload.Auth0ID == "") {
	// 	auth0Id := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
	// 	userPayload.Auth0ID = auth0Id
	// }

	user := c.store.CreateUser(userPayload)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func (c *Controllers) ModifyUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	err := json.NewDecoder(r.Body).Decode(&userPayload)

	if (err != nil) {
		w.WriteHeader(400)
		return
	}

	user := c.store.ModifyUser(userID, userPayload)

	render.JSON(w, r, user)
}

func (c *Controllers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	c.store.DeleteUser(userID)

	w.WriteHeader(http.StatusNoContent)
}
