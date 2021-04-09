package main

import (
  "log"
  "net/http"

  "go_server/internal/routes"
)

func main() {
	router := routes.Init()
	log.Fatal(http.ListenAndServe(":8080", router))
}
