package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	UsersService "go_server/internal/services/users"
)

func Get(w http.ResponseWriter, r *http.Request) {
	userID := uuid.Must(uuid.Parse(chi.URLParam(r, "userID")))

	user := UsersService.Get(userID)

	render.JSON(w, r, user)
}
