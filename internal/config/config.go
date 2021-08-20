package config

import (
	"fmt"
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
	RouterMaxAge   int
}

func NewConfig() (*Config, error) {
	dbSSL, dbSSLError := strconv.ParseBool(os.Getenv("DB_SSL"))

	if dbSSLError != nil {
		return nil, fmt.Errorf("error parsing dbSSL: %w", dbSSLError)
	}

	config := &Config{
		Port:           os.Getenv("PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBPort:         os.Getenv("DB_PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBName:         os.Getenv("DB_NAME"),
		DBSSL:          dbSSL,
		Auth0Domain:    os.Getenv("AUTH0_DOMAIN"),
		Auth0Audience:  os.Getenv("AUTH0_AUDIENCE"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		RouterMaxAge:   300,
	}

	return config, nil
}
