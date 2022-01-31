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

func (mockStore *MockBank) CreateToken(userID string) (string, error) {
	args := mockStore.Called(userID)

	return args.String(0), args.Error(1)
}

func (mockStore *MockBank) GetAccessToken(publicToken string) (string, error) {
	args := mockStore.Called(publicToken)

	return args.String(0), args.Error(1)
}

func (mockStore *MockBank) GetAccount(accessToken string) (string, error) {
	args := mockStore.Called(accessToken)

	return args.String(0), args.Error(1)
}
