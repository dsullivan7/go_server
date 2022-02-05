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

func (mockBroker *MockBroker) CreateAccount(
	givenName string,
	familyName string,
	dateOfBirth string,
	taxID string,
	emailAddress string,
	phoneNumber string,
	streetAddress string,
	city string,
	state string,
	postalCode string,
	fundingSource string,
	ipAddress string,
) (string, error) {
	args := mockBroker.Called(
		givenName,
		familyName,
		dateOfBirth,
		taxID,
		emailAddress,
		phoneNumber,
		streetAddress,
		city,
		state,
		postalCode,
		fundingSource,
		ipAddress,
	)

	return args.String(0), args.Error(1)
}
