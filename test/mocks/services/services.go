package services

import (
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func NewMockService() *MockService {
	return &MockService{}
}
