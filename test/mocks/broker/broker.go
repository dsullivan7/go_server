package broker

import (
	"github.com/stretchr/testify/mock"
)

type MockBroker struct {
	mock.Mock
}

func NewMockBroker() *MockBroker {
	return &MockBroker{}
}

func (mockBroker *MockBroker) CreateAccount(emailAddress string, phoneNumber string) (string, error) {
	args := mockBroker.Called(emailAddress, phoneNumber)

	return args.String(0), args.Error(1)
}
