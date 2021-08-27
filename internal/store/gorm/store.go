package gorm

import (
	"go_server/internal/store"

	"gorm.io/gorm"
)

type Store struct {
	database *gorm.DB
}

func NewStore(database *gorm.DB) store.Store {
	return &Store{database: database}
}
