package controllers

import (
	"net/http"
	"encoding/json"

	"go_server/internal/services"
	"go_server/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"

  "github.com/dgrijalva/jwt-go"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user := services.GetUser(userID)

	render.JSON(w, r, user)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	users := services.ListUsers(&models.User{})

	render.JSON(w, r, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	json.NewDecoder(r.Body).Decode(&userPayload)

	if (userPayload.Auth0ID == "") {
		auth0Id := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["sub"].(string)
		userPayload.Auth0ID = auth0Id
	}

	user := services.CreateUser(userPayload)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	json.NewDecoder(r.Body).Decode(&userPayload)

	user := services.ModifyUser(userID, userPayload)

	render.JSON(w, r, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	services.DeleteUser(userID)

	w.WriteHeader(http.StatusNoContent)
}
