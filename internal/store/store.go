package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	TruncateAll()

	GetUser(userID uuid.UUID) (*models.User, error)
	ListUsers(query map[string]interface{}) ([]models.User, error)
	CreateUser(userPayload models.User) (*models.User, error)
	ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) error

	GetReview(userID uuid.UUID) (*models.Review, error)
	ListReviews(query map[string]interface{}) ([]models.Review, error)
	CreateReview(userPayload models.Review) (*models.Review, error)
	ModifyReview(userID uuid.UUID, userPayload models.Review) (*models.Review, error)
	DeleteReview(userID uuid.UUID) error
}

type GormStore struct {
	database *gorm.DB
}

func NewGormStore(database *gorm.DB) Store {
	return &GormStore{database: database}
}

func (gormStore *GormStore) TruncateAll() {
	gormStore.database.Exec(`
		do $$
		begin
			execute (
				select 'truncate table ' || string_agg('"' || tablename || '"', ', ')
				from pg_tables
				where schemaname = 'public'
			);
		end;
		$$
	`)
}
