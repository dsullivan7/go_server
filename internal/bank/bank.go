package bank

type Bank interface {
	CreateToken(userID string) (string, error)
	GetAccessToken(publicToken string) (string, error)
	GetAccount(accessToken string) (string, error)
}
