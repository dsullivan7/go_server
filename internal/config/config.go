package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port string

	DBHost     string
	DBName     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBSSL      bool

	TwoCaptchaKey string

	Auth0Domain   string
	Auth0Audience string

	PlaidClientID    string
	PlaidSecret      string
	PlaidRedirectURI string

	AlpacaAPIKey    string
	AlpacaAPISecret string
	AlpacaAPIURL    string

	AllowedOrigins []string
	RouterMaxAge   int
}

func NewConfig() (*Config, error) {
	dbSSL, dbSSLError := strconv.ParseBool(os.Getenv("DB_SSL"))

	if dbSSLError != nil {
		return nil, fmt.Errorf("error parsing dbSSL: %w", dbSSLError)
	}

	config := &Config{
		Port:             os.Getenv("PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBPort:           os.Getenv("DB_PORT"),
		DBHost:           os.Getenv("DB_HOST"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSL:            dbSSL,
		TwoCaptchaKey:    os.Getenv("TWO_CAPTCHA_KEY"),
		PlaidClientID:    os.Getenv("PLAID_CLIENT_ID"),
		PlaidSecret:      os.Getenv("PLAID_SECRET"),
		PlaidRedirectURI: os.Getenv("PLAID_REDIRECT_URI"),
		AlpacaAPIKey:     os.Getenv("ALPACA_API_KEY"),
		AlpacaAPISecret:  os.Getenv("ALPACA_API_SECRET"),
		AlpacaAPIURL:     os.Getenv("ALPACA_API_URL"),
		Auth0Domain:      os.Getenv("AUTH0_DOMAIN"),
		Auth0Audience:    os.Getenv("AUTH0_AUDIENCE"),
		AllowedOrigins:   strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		RouterMaxAge:     300,
	}

	return config, nil
}
