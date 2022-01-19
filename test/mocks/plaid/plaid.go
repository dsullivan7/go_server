package plaid

import (
	"github.com/stretchr/testify/mock"
)

type MockPlaidClient struct {
	mock.Mock
}

func NewMockPlaidClient() *MockPlaidClient {
	return &MockPlaidClient{}
}

func (mockStore *MockPlaidClient) CreatePlaidToken(userID string) (string, error) {
	args := mockStore.Called(userID)

	return args.String(0), args.Error(1)
}
