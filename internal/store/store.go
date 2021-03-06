package store

import (
	"go_server/internal/models"

	"github.com/google/uuid"
)

type Store interface {
	GetUser(userID uuid.UUID) (*models.User, error)
	ListUsers(query map[string]interface{}) ([]models.User, error)
	CreateUser(userPayload models.User) (*models.User, error)
	ModifyUser(userID uuid.UUID, userPayload models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) error

	GetTag(tagID uuid.UUID) (*models.Tag, error)
	ListTags(query map[string]interface{}) ([]models.Tag, error)
	CreateTag(tagPayload models.Tag) (*models.Tag, error)
	ModifyTag(tagID uuid.UUID, tagPayload models.Tag) (*models.Tag, error)
	DeleteTag(tagID uuid.UUID) error

	GetBankAccount(bankAccountID uuid.UUID) (*models.BankAccount, error)
	ListBankAccounts(query map[string]interface{}) ([]models.BankAccount, error)
	CreateBankAccount(bankAccountPayload models.BankAccount) (*models.BankAccount, error)
	ModifyBankAccount(bankAccountID uuid.UUID, bankAccountPayload models.BankAccount) (*models.BankAccount, error)
	DeleteBankAccount(bankAccountID uuid.UUID) error

	GetBankTransfer(bankTransferID uuid.UUID) (*models.BankTransfer, error)
	ListBankTransfers(query map[string]interface{}) ([]models.BankTransfer, error)
	CreateBankTransfer(bankTransferPayload models.BankTransfer) (*models.BankTransfer, error)
	ModifyBankTransfer(bankTransferID uuid.UUID, bankTransferPayload models.BankTransfer) (*models.BankTransfer, error)
	DeleteBankTransfer(bankTransferID uuid.UUID) error

	GetBrokerageAccount(brokerageAccountID uuid.UUID) (*models.BrokerageAccount, error)
	ListBrokerageAccounts(query map[string]interface{}) ([]models.BrokerageAccount, error)
	CreateBrokerageAccount(brokerageAccountPayload models.BrokerageAccount) (*models.BrokerageAccount, error)
	ModifyBrokerageAccount(
		brokerageAccountID uuid.UUID,
		brokerageAccountPayload models.BrokerageAccount,
	) (*models.BrokerageAccount, error)
	DeleteBrokerageAccount(brokerageAccountID uuid.UUID) error

	GetPortfolio(portfolioID uuid.UUID) (*models.Portfolio, error)
	ListPortfolios(query map[string]interface{}) ([]models.Portfolio, error)
	CreatePortfolio(portfolioPayload models.Portfolio) (*models.Portfolio, error)
	ModifyPortfolio(portfolioID uuid.UUID, portfolioPayload models.Portfolio) (*models.Portfolio, error)
	DeletePortfolio(portfolioID uuid.UUID) error

	GetPortfolioTag(portfolioTagID uuid.UUID) (*models.PortfolioTag, error)
	ListPortfolioTags(query map[string]interface{}) ([]models.PortfolioTag, error)
	CreatePortfolioTag(portfolioTagPayload models.PortfolioTag) (*models.PortfolioTag, error)
	ModifyPortfolioTag(
		portfolioTagID uuid.UUID,
		portfolioTagPayload models.PortfolioTag,
	) (*models.PortfolioTag, error)
	DeletePortfolioTag(portfolioTagID uuid.UUID) error

	GetSecurity(securityID uuid.UUID) (*models.Security, error)
	ListSecurities(query map[string]interface{}) ([]models.Security, error)
	CreateSecurity(securityPayload models.Security) (*models.Security, error)
	ModifySecurity(securityID uuid.UUID, securityPayload models.Security) (*models.Security, error)
	DeleteSecurity(securityID uuid.UUID) error

	GetSecurityTag(securityTagID uuid.UUID) (*models.SecurityTag, error)
	ListSecurityTags(query map[string]interface{}) ([]models.SecurityTag, error)
	CreateSecurityTag(securityTagPayload models.SecurityTag) (*models.SecurityTag, error)
	ModifySecurityTag(
		securityTagID uuid.UUID,
		securityTagPayload models.SecurityTag,
	) (*models.SecurityTag, error)
	DeleteSecurityTag(securityTagID uuid.UUID) error

	GetOrder(orderID uuid.UUID) (*models.Order, error)
	ListOrders(query map[string]interface{}) ([]models.Order, error)
	CreateOrder(orderPayload models.Order) (*models.Order, error)
	ModifyOrder(orderID uuid.UUID, orderPayload models.Order) (*models.Order, error)
	DeleteOrder(orderID uuid.UUID) error

	GetItem(itemID uuid.UUID) (*models.Item, error)
	ListItems(query map[string]interface{}) ([]models.Item, error)
	CreateItem(itemPayload models.Item) (*models.Item, error)
	ModifyItem(itemID uuid.UUID, itemPayload models.Item) (*models.Item, error)
	DeleteItem(itemID uuid.UUID) error

	GetInvoice(itemID uuid.UUID) (*models.Invoice, error)
	ListInvoices(query map[string]interface{}) ([]models.Invoice, error)
	CreateInvoice(itemPayload models.Invoice) (*models.Invoice, error)
	ModifyInvoice(itemID uuid.UUID, itemPayload models.Invoice) (*models.Invoice, error)
	DeleteInvoice(itemID uuid.UUID) error

	CreateCredential(credentialPayload models.Credential) (*models.Credential, error)
}
