package main

import (
  "log"
  "net/http"

  "github.com/go-chi/chi"
  UserRoutes "go_server/internal/routes/users"
)

func main() {
	router := chi.NewRouter()
  router.Mount("/api/users", UserRoutes.Routes())
	log.Fatal(http.ListenAndServe(":8080", router))
}
