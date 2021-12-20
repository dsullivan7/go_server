package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (c *Controllers) GetHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "healthy"}

	render.JSON(w, r, response)
}
