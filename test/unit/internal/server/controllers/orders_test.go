package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	"go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestOrderGet(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	orderID := uuid.New()
	userID := uuid.New()
	portfolioID := uuid.New()
	alpacaOrderID := "alpacaOrderID"

	order := models.Order{
		OrderID:       orderID,
		UserID:        &userID,
		PortfolioID:   &portfolioID,
		AlpacaOrderID: &alpacaOrderID,
		Amount:        12345,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	testServer.Store.On("GetOrder", orderID).Return(&order, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_id", orderID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().GetOrder(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var orderResponse models.Order
	errDecoder := decoder.Decode(&orderResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, orderResponse.OrderID, order.OrderID)
	assert.Equal(t, orderResponse.UserID, order.UserID)
	assert.Equal(t, orderResponse.PortfolioID, order.PortfolioID)
	assert.Equal(t, orderResponse.Amount, order.Amount)
	assert.Equal(t, orderResponse.Side, order.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, order.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, order.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestOrderList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	orderID1 := uuid.New()
	userID1 := uuid.New()
	portfolioID1 := uuid.New()
	alpacaOrderID1 := "alpacaOrderID1"

	order1 := models.Order{
		OrderID:       orderID1,
		UserID:        &userID1,
		PortfolioID:   &portfolioID1,
		AlpacaOrderID: &alpacaOrderID1,
		Amount:        12345,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	orderID2 := uuid.New()
	userID2 := uuid.New()
	portfolioID2 := uuid.New()
	alpacaOrderID2 := "alpacaOrderID2"

	order2 := models.Order{
		OrderID:       orderID2,
		UserID:        &userID2,
		PortfolioID:   &portfolioID2,
		AlpacaOrderID: &alpacaOrderID2,
		Amount:        34567,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	testServer.Store.On("ListOrders", map[string]interface{}{}).Return([]models.Order{order1, order2}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/orders",
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListOrders(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var ordersFound []models.Order
	errDecoder := decoder.Decode(&ordersFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(ordersFound))

	var orderResponse models.Order

	for _, value := range ordersFound {
		if value.OrderID == order1.OrderID {
			orderResponse = value

			break
		}
	}

	assert.Equal(t, orderResponse.OrderID, order1.OrderID)
	assert.Equal(t, orderResponse.UserID, order1.UserID)
	assert.Equal(t, orderResponse.PortfolioID, order1.PortfolioID)
	assert.Equal(t, orderResponse.AlpacaOrderID, order1.AlpacaOrderID)
	assert.Equal(t, orderResponse.Amount, order1.Amount)
	assert.Equal(t, orderResponse.Side, order1.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, order1.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, order1.UpdatedAt, 0)

	for _, value := range ordersFound {
		if value.OrderID == order2.OrderID {
			orderResponse = value

			break
		}
	}

	assert.Equal(t, orderResponse.OrderID, order2.OrderID)
	assert.Equal(t, orderResponse.UserID, order2.UserID)
	assert.Equal(t, orderResponse.PortfolioID, order2.PortfolioID)
	assert.Equal(t, orderResponse.AlpacaOrderID, order2.AlpacaOrderID)
	assert.Equal(t, orderResponse.Amount, order2.Amount)
	assert.Equal(t, orderResponse.Side, order2.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, order2.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, order2.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestOrderListQueryParams(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	orderID := uuid.New()
	userID := uuid.New()
	portfolioID := uuid.New()
	alpacaOrderID := "alpacaOrderID"
	amount := 12345

	order := models.Order{
		OrderID:       orderID,
		UserID:        &userID,
		PortfolioID:   &portfolioID,
		AlpacaOrderID: &alpacaOrderID,
		Amount:        amount,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	testServer.Store.On(
		"ListOrders",
		map[string]interface{}{"user_id": userID.String()},
	).Return([]models.Order{order}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/api/orders?user_id=", userID),
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListOrders(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var ordersFound []models.Order
	errDecoder := decoder.Decode(&ordersFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 1, len(ordersFound))

	orderResponse := ordersFound[0]

	assert.Equal(t, orderResponse.OrderID, order.OrderID)
	assert.Equal(t, orderResponse.UserID, order.UserID)
	assert.Equal(t, orderResponse.PortfolioID, order.PortfolioID)
	assert.Equal(t, orderResponse.AlpacaOrderID, order.AlpacaOrderID)
	assert.Equal(t, orderResponse.Amount, order.Amount)
	assert.Equal(t, orderResponse.Side, order.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, order.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, order.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestOrderCreate(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()
	parentOrderID := uuid.New()
	amount := 12345

	jsonStr := []byte(fmt.Sprintf(
		`{
				"user_id": "%s",
				"amount": %d,
				"side": "buy"
			}`,
		userID.String(),
		amount,
	))

	orderPayloadParent := models.Order{
		UserID:      &userID,
		Amount:      amount,
		Side:        "buy",
	}

	orderCreatedParent := models.Order{
		OrderID:     parentOrderID,
		UserID:      &userID,
		Amount:      amount,
		Side:        "buy",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	orderPayloadChild := models.Order{
		UserID:        &userID,
		ParentOrderID: &parentOrderID,
		Amount:        amount,
		Side:          "buy",
	}

	orderCreatedChild := models.Order{
		OrderID:       uuid.New(),
		UserID:        &userID,
		Amount:        amount,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	testServer.Store.On("CreateOrder", orderPayloadParent).Return(&orderCreatedParent, nil)
	testServer.Store.On("CreateOrder", orderPayloadChild).Return(&orderCreatedChild, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/orders",
		bytes.NewBuffer(jsonStr),
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().CreateOrder(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var orderResponse models.Order
	errDecoder := decoder.Decode(&orderResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, orderResponse.OrderID, orderCreatedParent.OrderID)
	assert.Equal(t, orderResponse.UserID, orderCreatedParent.UserID)
	assert.Equal(t, orderResponse.Amount, orderCreatedParent.Amount)
	assert.Equal(t, orderResponse.Side, orderCreatedParent.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, orderCreatedParent.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, orderCreatedParent.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestOrderModify(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()
	portfolioID := uuid.New()
	alpacaOrderID := "alpacaOrderID"
	amount := 12345

	jsonStr := []byte(`{}`)

	orderPayload := models.Order{}

	orderModified := models.Order{
		OrderID:       uuid.New(),
		UserID:        &userID,
		PortfolioID:   &portfolioID,
		AlpacaOrderID: &alpacaOrderID,
		Amount:        amount,
		Side:          "buy",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	testServer.Store.On("ModifyOrder", orderModified.OrderID, orderPayload).Return(&orderModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/orders",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_id", orderModified.OrderID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ModifyOrder(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var orderResponse models.Order
	errDecoder := decoder.Decode(&orderResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, orderResponse.OrderID, orderModified.OrderID)
	assert.Equal(t, orderResponse.UserID, orderModified.UserID)
	assert.Equal(t, orderResponse.PortfolioID, orderModified.PortfolioID)
	assert.Equal(t, orderResponse.AlpacaOrderID, orderModified.AlpacaOrderID)
	assert.Equal(t, orderResponse.Amount, orderModified.Amount)
	assert.Equal(t, orderResponse.Side, orderModified.Side)
	assert.WithinDuration(t, orderResponse.CreatedAt, orderModified.CreatedAt, 0)
	assert.WithinDuration(t, orderResponse.UpdatedAt, orderModified.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestOrderDelete(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	orderID := uuid.New()

	testServer.Store.On("DeleteOrder", orderID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/orders",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_id", orderID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().DeleteOrder(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	testServer.Store.AssertExpectations(t)
}
