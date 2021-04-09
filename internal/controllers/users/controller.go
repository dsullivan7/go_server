package users

import (
  "net/http"

  "github.com/go-chi/chi"
  "github.com/go-chi/render"

  UsersService "go_server/internal/services/users"
)

func Get(w http.ResponseWriter, r *http.Request) {
    userId := chi.URLParam(r, "userId")

    user := UsersService.Get(userId)

    render.JSON(w, r, user)
}
