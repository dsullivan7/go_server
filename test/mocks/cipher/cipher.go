package broker

import (
	"github.com/stretchr/testify/mock"
)

type MockCipher struct {
	mock.Mock
}

func NewMockCipher() *MockCipher {
	return &MockCipher{}
}

func (mockCipher *MockCipher) Encrypt(
	phrase string,
	key string,
) (string, error) {
	args := mockCipher.Called(phrase, key)

	return args.String(0), args.Error(1)
}

func (mockCipher *MockCipher) Decrypt(
	phrase string,
	key string,
) (string, error) {
	args := mockCipher.Called(phrase, key)

	return args.String(0), args.Error(1)
}
