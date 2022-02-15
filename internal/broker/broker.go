package broker

type Account struct {
	Cash    float64
	AccountID string
}

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
	GetAccount(accountID string) (*Account, error)
	DeleteAccount(accountID string) error
	CreateOrder(
		accountID string,
		symbol string,
		quantity float64,
		side string,
	) (string, error)
	CreateTransfer(
		accountID string,
		relationshipID string,
		amount float64,
		direction string,
	) (string, error)
	CreateACHRelationship(
		accountID string,
		processorToken string,
	) (string, error)
}
