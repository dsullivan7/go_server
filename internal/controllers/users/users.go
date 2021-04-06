package users

import (
  "net/http"

  "github.com/go-chi/render"

  models "go_server/internal/models"
)

func Get(w http.ResponseWriter, r *http.Request) {
    user := models.User{
      FirstName: "blah",
    }
    render.JSON(w, r, user)
}
