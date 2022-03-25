package controllers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/broker"
	"go_server/internal/models"
	"go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestPositionsList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID := uuid.New()
	alpacaAccountID := "alpacaAccountID"

	brokerageAccount := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		AlpacaAccountID:    &alpacaAccountID,
	}

	positionID1 := "positionID1"

	position1 := broker.Position{
		PositionID:  positionID1,
		Symbol:      "symbol1",
		MarketValue: 123.45,
	}

	positionID2 := "positionID2"

	position2 := broker.Position{
		PositionID:  positionID2,
		Symbol:      "symbol2",
		MarketValue: 234.56,
	}

	testServer.Store.On("GetBrokerageAccount", brokerageAccountID).Return(&brokerageAccount, nil)
	testServer.Broker.On("ListPositions", alpacaAccountID).Return(
		[]broker.Position{position1, position2},
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/api/positions?brokerage_account_id=", brokerageAccountID),
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("brokerage_account_id", brokerageAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListPositions(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var positionsFound []broker.Position
	errDecoder := decoder.Decode(&positionsFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(positionsFound))

	positionResponse1 := positionsFound[0]

	assert.Equal(t, positionResponse1.PositionID, position1.PositionID)
	assert.Equal(t, positionResponse1.Symbol, position1.Symbol)
	assert.Equal(t, positionResponse1.MarketValue, position1.MarketValue)

	positionResponse2 := positionsFound[1]

	assert.Equal(t, positionResponse2.PositionID, position2.PositionID)
	assert.Equal(t, positionResponse2.Symbol, position2.Symbol)
	assert.Equal(t, positionResponse2.MarketValue, position2.MarketValue)

	testServer.Broker.AssertExpectations(t)
	testServer.Store.AssertExpectations(t)
}
