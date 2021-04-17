package db

import (
  "go_server/internal/models"

  "fmt"

  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

var DB *gorm.DB

func Connect() {
  dbName := "go_server"
  dbHost := "localhost"

  DSN := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbHost, dbName)

  fmt.Println(DSN)

  database, err := gorm.Open(
    postgres.New(
      postgres.Config{
        DSN: DSN,
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
