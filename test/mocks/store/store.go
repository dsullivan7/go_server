package store

import (
	"go_server/internal/models"
	"go_server/internal/store"

	"github.com/google/uuid"
)

type MockStore struct{}

func NewMockStore() store.Store {
	return &MockStore{}
}

// User.
func (mockStore *MockStore) GetUser(userID uuid.UUID) (*models.User, error) {
	return &models.User{}, nil
}

func (mockStore *MockStore) ListUsers(query map[string]interface{}) ([]models.User, error) {
	return []models.User{models.User{}}, nil
}

func (mockStore *MockStore) CreateUser(userPayload models.User) (*models.User, error) {
	return &userPayload, nil
}

func (mockStore *MockStore) ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error) {
	return &userPayload, nil
}

func (mockStore *MockStore) DeleteUser(userID uuid.UUID) error {
	return nil
}

// Review.
func (mockStore *MockStore) GetReview(industryID uuid.UUID) (*models.Review, error) {
	return &models.Review{}, nil
}

func (mockStore *MockStore) ListReviews(query map[string]interface{}) ([]models.Review, error) {
	return []models.Review{models.Review{}}, nil
}

func (mockStore *MockStore) CreateReview(reviewPayload models.Review) (*models.Review, error) {
	return &reviewPayload, nil
}

func (mockStore *MockStore) ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) (*models.Review, error) {
	return &reviewPayload, nil
}

func (mockStore *MockStore) DeleteReview(reviewID uuid.UUID) error {
	return nil
}

// Industry.
func (mockStore *MockStore) GetIndustry(industryID uuid.UUID) (*models.Industry, error) {
	return &models.Industry{}, nil
}

func (mockStore *MockStore) ListIndustries(query map[string]interface{}) ([]models.Industry, error) {
	return []models.Industry{models.Industry{}}, nil
}

func (mockStore *MockStore) CreateIndustry(industryPayload models.Industry) (*models.Industry, error) {
	return &industryPayload, nil
}

func (mockStore *MockStore) ModifyIndustry(
	industryID uuid.UUID,
	industryPayload models.Industry,
) (*models.Industry, error) {
	return &industryPayload, nil
}

func (mockStore *MockStore) DeleteIndustry(industryID uuid.UUID) error {
	return nil
}

// UserIndustry.
func (mockStore *MockStore) GetUserIndustry(userIndustryID uuid.UUID) (*models.UserIndustry, error) {
	return &models.UserIndustry{}, nil
}

func (mockStore *MockStore) ListUserIndustries(query map[string]interface{}) ([]models.UserIndustry, error) {
	return []models.UserIndustry{models.UserIndustry{}}, nil
}

func (mockStore *MockStore) CreateUserIndustry(userIndustryPayload models.UserIndustry) (*models.UserIndustry, error) {
	return &userIndustryPayload, nil
}

func (mockStore *MockStore) ModifyUserIndustry(
	userIndustryID uuid.UUID,
	userIndustryPayload models.UserIndustry,
) (*models.UserIndustry, error) {
	return &userIndustryPayload, nil
}

func (mockStore *MockStore) DeleteUserIndustry(userIndustryID uuid.UUID) error {
	return nil
}
