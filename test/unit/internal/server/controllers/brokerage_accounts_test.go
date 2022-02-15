package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/broker"
	"go_server/internal/models"
	"go_server/internal/server/controllers"
	"go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestBrokerageAccountGet(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID := uuid.New()
	userID := uuid.New()

	alpacaAccountID := "alpacaAccountID"

	brokerageAccount := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		UserID:             &userID,
		AlpacaAccountID:    &alpacaAccountID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	testServer.Store.On("GetBrokerageAccount", brokerageAccountID).Return(&brokerageAccount, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/brokerage-accounts",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("brokerage_account_id", brokerageAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().GetBrokerageAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var brokerageAccountResponse models.BrokerageAccount
	errDecoder := decoder.Decode(&brokerageAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccount.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccount.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccount.AlpacaAccountID)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccount.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccount.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBrokerageAccountList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID1 := uuid.New()
	userID1 := uuid.New()

	alpacaAccountID1 := "alpacaAccountID1"
	cash1 := 234.56

	brokerageAccount1 := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID1,
		UserID:             &userID1,
		AlpacaAccountID:    &alpacaAccountID1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	alpacaAccount1 := broker.Account{
		AccountID: alpacaAccountID1,
		Cash:      cash1,
	}

	brokerageAccountID2 := uuid.New()
	userID2 := uuid.New()

	alpacaAccountID2 := "alpacaAccountID2"
	cash2 := 456.56

	brokerageAccount2 := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID2,
		UserID:             &userID2,
		AlpacaAccountID:    &alpacaAccountID2,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	alpacaAccount2 := broker.Account{
		AccountID: alpacaAccountID2,
		Cash:      cash2,
	}

	testServer.Store.On(
		"ListBrokerageAccounts",
		map[string]interface{}{},
	).Return([]models.BrokerageAccount{brokerageAccount1, brokerageAccount2}, nil)
	testServer.Broker.On("GetAccount", alpacaAccountID1).Return(&alpacaAccount1, nil)
	testServer.Broker.On("GetAccount", alpacaAccountID2).Return(&alpacaAccount2, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/brokerage-accounts",
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBrokerageAccounts(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var brokerageAccountsFound []controllers.BrokerageAccountResponse
	errDecoder := decoder.Decode(&brokerageAccountsFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(brokerageAccountsFound))

	var brokerageAccountResponse controllers.BrokerageAccountResponse

	for _, value := range brokerageAccountsFound {
		if value.BrokerageAccountID == brokerageAccount1.BrokerageAccountID {
			brokerageAccountResponse = value

			break
		}
	}

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccount1.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccount1.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccount1.AlpacaAccountID)
	assert.Equal(t, *brokerageAccountResponse.Cash, cash1)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccount1.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccount1.UpdatedAt, 0)

	for _, value := range brokerageAccountsFound {
		if value.BrokerageAccountID == brokerageAccount2.BrokerageAccountID {
			brokerageAccountResponse = value

			break
		}
	}

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccount2.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccount2.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccount2.AlpacaAccountID)
	assert.Equal(t, *brokerageAccountResponse.Cash, cash2)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccount2.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccount2.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBrokerageAccountListQueryParams(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()
	brokerageAccountID := uuid.New()

	alpacaAccountID := "alpacaAccountID"
	cash := 123.45

	brokerageAccount := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		UserID:             &userID,
		AlpacaAccountID:    &alpacaAccountID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	alpacaAccount := broker.Account{
		AccountID: alpacaAccountID,
		Cash:      cash,
	}

	testServer.Store.
		On("ListBrokerageAccounts", map[string]interface{}{"user_id": userID.String()}).
		Return([]models.BrokerageAccount{brokerageAccount}, nil)
	testServer.Broker.On("GetAccount", alpacaAccountID).Return(&alpacaAccount, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/api/brokerage-accounts?user_id=", userID),
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBrokerageAccounts(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var brokerageAccountsFound []controllers.BrokerageAccountResponse
	errDecoder := decoder.Decode(&brokerageAccountsFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 1, len(brokerageAccountsFound))

	brokerageAccountResponse := brokerageAccountsFound[0]

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccount.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccount.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccount.AlpacaAccountID)
	assert.Equal(t, *brokerageAccountResponse.Cash, cash)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccount.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccount.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBrokerageAccountCreate(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID := uuid.New()
	userID := uuid.New()

	emailAddress := "emailAddress"
	firstName := "firstName"
	lastName := "lastName"
	dateOfBirth := "dateOfBirth"
	taxID := "taxID"
	phoneNumber := "phoneNumber"
	streetAddress := "streetAddress"
	city := "city"
	state := "state"
	postalCode := "postalCode"
	fundingSource := "fundingSource"
	ipAddress := "ipAddress"

	alpacaAccountID := "alpacaAccountID"

	brokerageAccountPayload := models.BrokerageAccount{
		UserID:          &userID,
		AlpacaAccountID: &alpacaAccountID,
	}

	brokerageAccountCreated := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		UserID:             &userID,
		AlpacaAccountID:    &alpacaAccountID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	jsonStr := []byte(fmt.Sprintf(
		`{
				"user_id": "%s",
				"email_address": "%s",
				"first_name": "%s",
				"last_name": "%s",
				"tax_id": "%s",
				"date_of_birth": "%s",
				"phone_number": "%s",
				"street_address": "%s",
				"city": "%s",
				"state": "%s",
				"postal_code": "%s",
				"funding_source": "%s"
			}`,
		userID.String(),
		emailAddress,
		firstName,
		lastName,
		taxID,
		dateOfBirth,
		phoneNumber,
		streetAddress,
		city,
		state,
		postalCode,
		fundingSource,
	))

	testServer.Broker.
		On(
			"CreateAccount",
			firstName,
			lastName,
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
		).
		Return(alpacaAccountID, nil)
	testServer.Store.On("CreateBrokerageAccount", brokerageAccountPayload).Return(&brokerageAccountCreated, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/brokerage-accounts",
		bytes.NewBuffer(jsonStr),
	)

	req.RemoteAddr = ipAddress

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().CreateBrokerageAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var brokerageAccountResponse models.BrokerageAccount
	errDecoder := decoder.Decode(&brokerageAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccountCreated.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccountCreated.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccountCreated.AlpacaAccountID)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccountCreated.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccountCreated.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
	testServer.Broker.AssertExpectations(t)
}

func TestBrokerageAccountModify(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID := uuid.New()
	userID := uuid.New()

	alpacaAccountID := "alpacaAccountID"

	jsonStr := []byte(`{}`)

	brokerageAccountPayload := models.BrokerageAccount{}

	brokerageAccountModified := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		UserID:             &userID,
		AlpacaAccountID:    &alpacaAccountID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	testServer.Store.On(
		"ModifyBrokerageAccount",
		brokerageAccountModified.BrokerageAccountID,
		brokerageAccountPayload,
	).Return(&brokerageAccountModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/brokerage-accounts",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("brokerage_account_id", brokerageAccountModified.BrokerageAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ModifyBrokerageAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var brokerageAccountResponse models.BrokerageAccount
	errDecoder := decoder.Decode(&brokerageAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, brokerageAccountResponse.BrokerageAccountID, brokerageAccountModified.BrokerageAccountID)
	assert.Equal(t, *brokerageAccountResponse.UserID, *brokerageAccountModified.UserID)
	assert.Equal(t, *brokerageAccountResponse.AlpacaAccountID, *brokerageAccountModified.AlpacaAccountID)
	assert.WithinDuration(t, brokerageAccountResponse.CreatedAt, brokerageAccountModified.CreatedAt, 0)
	assert.WithinDuration(t, brokerageAccountResponse.UpdatedAt, brokerageAccountModified.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBrokerageAccountDelete(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	brokerageAccountID := uuid.New()

	testServer.Store.On("DeleteBrokerageAccount", brokerageAccountID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/brokerage-accounts",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("brokerage_account_id", brokerageAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().DeleteBrokerageAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	testServer.Store.AssertExpectations(t)
}
