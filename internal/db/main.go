package db

import (
  "go_server/internal/models"
  "go_server/internal/config"

  "fmt"
  "strings"

  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

var DB *gorm.DB

func Connect() {
  var DSN strings.Builder

  DSN.WriteString(fmt.Sprintf("host=%s dbname=%s", config.DBHost, config.DBName))

  if config.DBPort != "" {
    DSN.WriteString(fmt.Sprintf(" port=%s", config.DBPort))
  }

  if config.DBUser != "" {
    DSN.WriteString(fmt.Sprintf(" user=%s", config.DBUser))
  }

  if config.DBPassword != "" {
    DSN.WriteString(fmt.Sprintf(" password=%s", config.DBPassword))
  }

  if config.DBSSL != true {
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

  if err != nil {
    panic("Failed to connect to database!")
  }

  DB = database
}

func Migrate() {
  DB.Exec("create extension if not exists \"uuid-ossp\"")
  DB.AutoMigrate(&models.User{})
}
