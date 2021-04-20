package config

import (
	"os"
  "strconv"
)

var Port = os.Getenv("PORT")

var DBPort = os.Getenv("DB_PORT")
var DBName = os.Getenv("DB_NAME")
var DBHost = os.Getenv("DB_HOST")
var DBUser = os.Getenv("DB_USER")
var DBSSL, _ = strconv.ParseBool(os.Getenv("DB_SSL"))
var DBPassword = os.Getenv("DB_PASSWORD")
