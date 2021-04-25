package users

import (
	"net/http"
	"encoding/json"

	UsersService "go_server/internal/services/users"
	"go_server/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func Get(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user := UsersService.Get(userID)

	render.JSON(w, r, user)
}

func List(w http.ResponseWriter, r *http.Request) {
	users := UsersService.List()

	render.JSON(w, r, users)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User

	json.NewDecoder(r.Body).Decode(&userPayload)

	user := UsersService.Create(userPayload)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, user)
}

func Modify(w http.ResponseWriter, r *http.Request) {
	var userPayload models.User
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	json.NewDecoder(r.Body).Decode(&userPayload)

	user := UsersService.Modify(userID, userPayload)

	render.JSON(w, r, user)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	UsersService.Delete(userID)

	w.WriteHeader(http.StatusNoContent)
}
