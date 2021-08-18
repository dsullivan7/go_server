package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port string

	DBHost     string
	DBName     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBSSL      bool

	Auth0Domain   string
	Auth0Audience string

	AllowedOrigins string
}

func NewConfig() *Config {
	dbSSL, dbSSLError := strconv.ParseBool(os.Getenv("DB_SSL"))

	if dbSSLError != nil {
		os.Exit(1)
	}

	config := &Config{
		Port: os.Getenv("PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort: os.Getenv("DB_PORT"),
		DBHost: os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
		DBSSL: dbSSL,
		Auth0Domain: os.Getenv("AUTH0_DOMAIN"),
		Auth0Audience: os.Getenv("AUTH0_AUDIENCE"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}

	return config
}
