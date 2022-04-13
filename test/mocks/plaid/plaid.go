package plaid

import (
	"github.com/stretchr/testify/mock"
)

type MockPlaid struct {
	mock.Mock
}

func NewMockPlaid() *MockPlaid {
	return &MockPlaid{}
}

func (mockPlaid *MockPlaid) CreateToken(userID string) (string, error) {
	args := mockPlaid.Called(userID)

	return args.String(0), args.Error(1)
}

func (mockPlaid *MockPlaid) CreateProcessorToken(
	accessToken string,
	accountID string,
	accessor string,
) (string, error) {
	args := mockPlaid.Called(accessToken, accountID, accessor)

	return args.String(0), args.Error(1)
}

func (mockPlaid *MockPlaid) GetAccessToken(publicToken string) (string, error) {
	args := mockPlaid.Called(publicToken)

	return args.String(0), args.Error(1)
}

func (mockPlaid *MockPlaid) GetAccount(accessToken string) (string, string, error) {
	args := mockPlaid.Called(accessToken)

	return args.String(0), args.String(1), args.Error(2)
}

func (mockPlaid *MockPlaid) CreateTransferAuthorization(
	accessToken string,
	accountID string,
	originationAccountID string,
	amount string,
	transferType string,
	legalName string,
) (string, error) {
	args := mockPlaid.Called(accessToken, accountID, originationAccountID, amount, transferType, legalName)

	return args.String(0), args.Error(1)
}

func (mockPlaid *MockPlaid) CreateTransfer(
	accessToken string,
	accountID string,
	originationAccountID string,
	authorizationID string,
	amount string,
	transferType string,
	legalName string,
) (string, error) {
	args := mockPlaid.Called(accessToken, accountID, originationAccountID, authorizationID, amount, transferType, legalName)

	return args.String(0), args.Error(1)
}
