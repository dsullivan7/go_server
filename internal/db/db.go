package db

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	dbHost string,
	dbName string,
	dbPort string,
	dbUser string,
	dbPassword string,
	dbSSL bool,
) (*gorm.DB, error) {
	var DSN strings.Builder

	DSN.WriteString(fmt.Sprintf("host=%s dbname=%s", dbHost, dbName))

	if dbPort != "" {
		DSN.WriteString(fmt.Sprintf(" port=%s", dbPort))
	}

	if dbUser != "" {
		DSN.WriteString(fmt.Sprintf(" user=%s", dbUser))
	}

	if dbPassword != "" {
		DSN.WriteString(fmt.Sprintf(" password=%s", dbPassword))
	}

	if !dbSSL {
		DSN.WriteString(" sslmode=disable")
	}

	database, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: DSN.String(),
			},
		),
		&gorm.Config{},
	)

	return database, fmt.Errorf("failed to open db connection: %w", err)
}
