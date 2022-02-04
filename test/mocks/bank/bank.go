package bank

import (
	"github.com/stretchr/testify/mock"
)

type MockBank struct {
	mock.Mock
}

func NewMockBank() *MockBank {
	return &MockBank{}
}

func (mockBank *MockBank) CreateToken(userID string) (string, error) {
	args := mockBank.Called(userID)

	return args.String(0), args.Error(1)
}

func (mockBank *MockBank) GetAccessToken(publicToken string) (string, error) {
	args := mockBank.Called(publicToken)

	return args.String(0), args.Error(1)
}

func (mockBank *MockBank) GetAccount(accessToken string) (string, error) {
	args := mockBank.Called(accessToken)

	return args.String(0), args.Error(1)
}
