package broker

type Broker interface {
	CreateAccount(emailAddress string, phoneNumber string) (string, error)
}
