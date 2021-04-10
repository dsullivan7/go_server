package users

import (
  "log"
  "net/http"

  "github.com/go-chi/chi"
  "github.com/go-chi/render"

  "github.com/satori/go.uuid"

  UsersService "go_server/internal/services/users"
)

func Get(w http.ResponseWriter, r *http.Request) {
    userID, err := uuid.FromString(chi.URLParam(r, "userID"))

    if err != nil {
  		log.Printf("Something went wrong: %s", err)
  		return
  	}

    user := UsersService.Get(userID)

    render.JSON(w, r, user)
}
