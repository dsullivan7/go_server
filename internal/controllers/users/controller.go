package users

import (
	"net/http"

	UsersService "go_server/internal/services/users"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func Get(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user := UsersService.Get(userID)

	render.JSON(w, r, user)
}

func Create(w http.ResponseWriter, r *http.Request) {
	user := UsersService.Create()

	render.JSON(w, r, user)
}
