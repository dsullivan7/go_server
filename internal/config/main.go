package config

import (
	"os"
  "strconv"
)

var Port string

var AllowedOrigins string

var DBPort string
var DBName string
var DBHost string
var DBUser string
var DBSSL bool
var DBPassword string

var Auth0Domain string
var Auth0Audience string

func init() {
	Port = os.Getenv("PORT")

	AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")

	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBSSL, _ = strconv.ParseBool(os.Getenv("DB_SSL"))
	DBPassword = os.Getenv("DB_PASSWORD")

	Auth0Domain = os.Getenv("AUTH0_DOMAIN")
	Auth0Audience = os.Getenv("AUTH0_AUDIENCE")
}
