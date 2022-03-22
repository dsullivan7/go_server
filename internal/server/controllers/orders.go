package controllers

import (
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (c *Controllers) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	order, err := c.store.GetOrder(orderID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, order)
}

func (c *Controllers) ListOrders(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	orders, err := c.store.ListOrders(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, orders)
}

//nolint:funlen,cyclop
func (c *Controllers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&orderReq)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	portfolioID := uuid.Must(uuid.Parse(orderReq["portfolio_id"].(string)))
	userID := uuid.Must(uuid.Parse(orderReq["user_id"].(string)))

	portfolio, errPortfolio := c.store.GetPortfolio(portfolioID)

	if errPortfolio != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPortfolio})

		return
	}

	securities, errSecurities := c.store.ListSecurities(map[string]interface{}{})

	if errSecurities != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurities})

		return
	}

	securityTags, errSecurityTags := c.store.ListSecurityTags(map[string]interface{}{})

	if errSecurityTags != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errSecurityTags})

		return
	}

	portfolioTags, errPortfolioTags := c.store.ListPortfolioTags(
		map[string]interface{}{"portfolio_id": portfolioID.String()},
	)

	if errPortfolioTags != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errPortfolioTags})

		return
	}

	portfolioHoldings := c.services.ListPortfolioRecommendations(
		*portfolio,
		portfolioTags,
		securities,
		securityTags,
	)

	brokerageAccounts, errBrokerageAccounts := c.store.ListBrokerageAccounts(
		map[string]interface{}{"user_id": userID.String()},
	)

	if errBrokerageAccounts != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errBrokerageAccounts})

		return
	}

	alpacaAccountID := brokerageAccounts[0].AlpacaAccountID

	orderPayload := models.Order{
		UserID:      &userID,
		PortfolioID: &portfolioID,
		Amount:      orderReq["amount"].(float64),
		Side:        "buy",
	}

	order, err := c.store.CreateOrder(orderPayload)

	for _, portfolioHolding := range portfolioHoldings {
		alpacaOrderID, errOrder := c.broker.CreateOrder(
			*alpacaAccountID,
			portfolioHolding.Symbol,
			portfolioHolding.Amount*orderReq["amount"].(float64),
			"buy",
		)

		if errOrder != nil {
			c.utils.HandleError(w, r, errors.HTTPUserError{Err: errOrder})

			return
		}

		orderPayloadChild := models.Order{
			UserID:        &userID,
			ParentOrderID: &order.OrderID,
			AlpacaOrderID: &alpacaOrderID,
			Amount:        portfolioHolding.Amount * orderReq["amount"].(float64),
			Symbol:        &portfolioHolding.Symbol,
			Side:          "buy",
		}

		_, err := c.store.CreateOrder(orderPayloadChild)

		if err != nil {
			c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

			return
		}
	}

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, order)
}

func (c *Controllers) ModifyOrder(w http.ResponseWriter, r *http.Request) {
	var orderPayload models.Order

	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&orderPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	order, err := c.store.ModifyOrder(orderID, orderPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, order)
}

func (c *Controllers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := uuid.Must(uuid.Parse(chi.URLParam(r, "order_id")))

	err := c.store.DeleteOrder(orderID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
