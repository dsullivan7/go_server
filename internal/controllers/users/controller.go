package users

import (
  "net/http"

  "github.com/go-chi/chi"
  "github.com/go-chi/render"

  UsersService "go_server/internal/services/users"
)

func Get(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")

    user := UsersService.Get(userID)

    render.JSON(w, r, user)
}
