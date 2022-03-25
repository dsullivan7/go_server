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
) (string) {
	args := mockCipher.Called(phrase)

	return args.String(0)
}

func (mockCipher *MockCipher) Decrypt(
	phrase string,
) (string, error) {
	args := mockCipher.Called(phrase)

	return args.String(0), args.Error(1)
}
