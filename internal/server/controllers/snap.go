package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (c *Controllers) GetSnap(w http.ResponseWriter, r *http.Request) {
	text := c.crawler.Login("https://www.connectebt.com/nyebtclient/siteLogonClient.recip", "username", "password")

	response := map[string]string{"text": text}

	render.JSON(w, r, response)
}
