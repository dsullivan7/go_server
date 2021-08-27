package db

import (
	"fmt"
	"strings"

	"database/sql"

	// import the postgres driver.
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(sqlDB *sql.DB) (*gorm.DB, error) {
	database, err := gorm.Open(
		postgres.New(postgres.Config{
			Conn: sqlDB,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	return database, nil
}

func NewSQLConnection(
	dbHost string,
	dbName string,
	dbPort string,
	dbUser string,
	dbPassword string,
	dbSSL bool,
) (*sql.DB, error) {
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

	sqlDB, err := sql.Open("postgres", DSN.String())

	if err != nil {
		return nil, fmt.Errorf("failed to open sql connection: %w", err)
	}

	return sqlDB, nil
}
