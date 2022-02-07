package broker

type Broker interface {
	CreateAccount(
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
	) (string, error)
	DeleteAccount(accountID string) (error)
	CreateOrder(
		accountID string,
		symbol string,
		quantity float32,
		side string,
	) (string, error)
}
