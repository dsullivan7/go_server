package main

import (
	"go_server/internal/routes"
	"log"
	"net/http"
)

func main() {
	router := routes.Init()

	log.Println("Server started")

	log.Fatal(http.ListenAndServe(":8080", router))
}
