package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

type Store interface {
	GetUser(userID uuid.UUID) (*models.User, error)
	ListUsers(query map[string]interface{}) ([]models.User, error)
	CreateUser(userPayload models.User) (*models.User, error)
	ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) error

	GetReview(industryID uuid.UUID) (*models.Review, error)
	ListReviews(query map[string]interface{}) ([]models.Review, error)
	CreateReview(reviewPayload models.Review) (*models.Review, error)
	ModifyReview(industryID uuid.UUID, industryPayload models.Review) (*models.Review, error)
	DeleteReview(industryID uuid.UUID) error

	GetIndustry(industryID uuid.UUID) (*models.Industry, error)
	ListIndustries(query map[string]interface{}) ([]models.Industry, error)
	CreateIndustry(industryPayload models.Industry) (*models.Industry, error)
	ModifyIndustry(industryID uuid.UUID, industryPayload models.Industry) (*models.Industry, error)
	DeleteIndustry(industryID uuid.UUID) error

	GetUserIndustry(userIndustryID uuid.UUID) (*models.UserIndustry, error)
	ListUserIndustries(query map[string]interface{}) ([]models.UserIndustry, error)
	CreateUserIndustry(userIndustryPayload models.UserIndustry) (*models.UserIndustry, error)
	ModifyUserIndustry(userIndustryID uuid.UUID, userIndustryPayload models.UserIndustry) (*models.UserIndustry, error)
	DeleteUserIndustry(userIndustryID uuid.UUID) error
}
