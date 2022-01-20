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

	GetReview(reviewID uuid.UUID) (*models.Review, error)
	ListReviews(query map[string]interface{}) ([]models.Review, error)
	CreateReview(reviewPayload models.Review) (*models.Review, error)
	ModifyReview(reviewID uuid.UUID, reviewPayload models.Review) (*models.Review, error)
	DeleteReview(reviewID uuid.UUID) error

	GetIndustry(industryID uuid.UUID) (*models.Industry, error)
	ListIndustries(query map[string]interface{}) ([]models.Industry, error)
	CreateIndustry(industryPayload models.Industry) (*models.Industry, error)
	ModifyIndustry(industryID uuid.UUID, industryPayload models.Industry) (*models.Industry, error)
	DeleteIndustry(industryID uuid.UUID) error

	GetBankAccount(portfolioID uuid.UUID) (*models.BankAccount, error)
	ListBankAccounts(query map[string]interface{}) ([]models.BankAccount, error)
	CreateBankAccount(portfolioPayload models.BankAccount) (*models.BankAccount, error)
	ModifyBankAccount(portfolioID uuid.UUID, portfolioPayload models.BankAccount) (*models.BankAccount, error)
	DeleteBankAccount(portfolioID uuid.UUID) error

	GetPortfolio(portfolioID uuid.UUID) (*models.Portfolio, error)
	ListPortfolios(query map[string]interface{}) ([]models.Portfolio, error)
	CreatePortfolio(portfolioPayload models.Portfolio) (*models.Portfolio, error)
	ModifyPortfolio(portfolioID uuid.UUID, portfolioPayload models.Portfolio) (*models.Portfolio, error)
	DeletePortfolio(portfolioID uuid.UUID) error

	GetPortfolioIndustry(portfolioIndustryID uuid.UUID) (*models.PortfolioIndustry, error)
	ListPortfolioIndustries(query map[string]interface{}) ([]models.PortfolioIndustry, error)
	CreatePortfolioIndustry(portfolioIndustryPayload models.PortfolioIndustry) (*models.PortfolioIndustry, error)
	ModifyPortfolioIndustry(
		portfolioIndustryID uuid.UUID,
		portfolioIndustryPayload models.PortfolioIndustry,
	) (*models.PortfolioIndustry, error)
	DeletePortfolioIndustry(portfolioIndustryID uuid.UUID) error
}
