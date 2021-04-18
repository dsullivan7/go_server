package app

import (
	"go_server/internal/routes"
	"go_server/internal/db"
	"log"
	"net/http"
)

func Init() {
  db.Connect()

  router := routes.Init()

  log.Println("Server started")

  log.Fatal(http.ListenAndServe(":8080", router))
}
