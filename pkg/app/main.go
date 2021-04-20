package app

import (
	"go_server/internal/routes"
	"go_server/internal/db"
	"go_server/internal/config"
	"fmt"
	"log"
	"net/http"
)

func Init() {
  db.Connect()

  router := routes.Init()

  log.Println("Server started")

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router))
}