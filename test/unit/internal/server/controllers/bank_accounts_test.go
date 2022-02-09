package controllers_test

import (
	"time"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	"go_server/test/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestBankAccountGet(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankAccountID := uuid.New()
	userID := uuid.New()

	name := "name"
	plaidAccessToken := "plaidAccessToken"
	plaidAccountID := "plaidAccountID"
	alpacaACHRelationshipID := "alpacaACHRelationshipID"

	bankAccount := models.BankAccount{
		BankAccountID: bankAccountID,
		UserID:   &userID,
		Name: &name,
		PlaidAccessToken: &plaidAccessToken,
		PlaidAccountID: &plaidAccountID,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("GetBankAccount", bankAccountID).Return(&bankAccount, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/bank-accounts",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_account_id", bankAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().GetBankAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankAccountResponse models.BankAccount
	errDecoder := decoder.Decode(&bankAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccount.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccount.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccount.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccount.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccount.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccount.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccount.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccount.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankAccountList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankAccountID1 := uuid.New()
	userID1 := uuid.New()

	name1 := "name1"
	plaidAccessToken1 := "plaidAccessToken1"
	plaidAccountID1 := "plaidAccountID1"
	alpacaACHRelationshipID1 := "alpacaACHRelationshipID1"

	bankAccount1 := models.BankAccount{
		BankAccountID: bankAccountID1,
		UserID:   &userID1,
		Name: &name1,
		PlaidAccessToken: &plaidAccessToken1,
		PlaidAccountID: &plaidAccountID1,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bankAccountID2 := uuid.New()
	userID2 := uuid.New()

	name2 := "name2"
	plaidAccessToken2 := "plaidAccessToken2"
	plaidAccountID2 := "plaidAccountID2"
	alpacaACHRelationshipID2 := "alpacaACHRelationshipID2"

	bankAccount2 := models.BankAccount{
		BankAccountID: bankAccountID2,
		UserID:   &userID2,
		Name: &name2,
		PlaidAccessToken: &plaidAccessToken2,
		PlaidAccountID: &plaidAccountID2,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ListBankAccounts", map[string]interface{}{}).Return([]models.BankAccount{bankAccount1, bankAccount2}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/bank-accounts",
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBankAccounts(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankAccountsFound []models.BankAccount
	errDecoder := decoder.Decode(&bankAccountsFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(bankAccountsFound))

	var bankAccountResponse models.BankAccount

	for _, value := range bankAccountsFound {
		if value.BankAccountID == bankAccount1.BankAccountID {
			bankAccountResponse = value

			break
		}
	}

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccount1.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccount1.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccount1.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccount1.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccount1.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccount1.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccount1.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccount1.UpdatedAt, 0)

	for _, value := range bankAccountsFound {
		if value.BankAccountID == bankAccount2.BankAccountID {
			bankAccountResponse = value

			break
		}
	}

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccount2.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccount2.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccount2.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccount2.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccount2.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccount2.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccount2.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccount2.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankAccountListQueryParams(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()
	bankAccountID := uuid.New()

	name := "name"
	plaidAccessToken := "plaidAccessToken"
	plaidAccountID := "plaidAccountID"
	alpacaACHRelationshipID := "alpacaACHRelationshipID"

	bankAccount := models.BankAccount{
		BankAccountID: bankAccountID,
		UserID:   &userID,
		Name: &name,
		PlaidAccessToken: &plaidAccessToken,
		PlaidAccountID: &plaidAccountID,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ListBankAccounts", map[string]interface{}{ "user_id": userID.String() }).Return([]models.BankAccount{bankAccount}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/api/bank-accounts?user_id=", userID),
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBankAccounts(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankAccountsFound []models.BankAccount
	errDecoder := decoder.Decode(&bankAccountsFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 1, len(bankAccountsFound))

	bankAccountResponse := bankAccountsFound[0]

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccount.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccount.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccount.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccount.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccount.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccount.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccount.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccount.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankAccountCreate(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankAccountID := uuid.New()
	userID := uuid.New()

	name := "name"
	plaidAccessToken := "plaidAccessToken"
	plaidPublicToken := "plaidPublicToken"
	plaidProcessorToken := "plaidProcessorToken"
	plaidAccountID := "plaidAccountID"
	alpacaACHRelationshipID := "alpacaACHRelationshipID"

	bankAccountPayload := models.BankAccount{
		UserID:           &userID,
		Name:             &name,
		PlaidAccountID:   &plaidAccountID,
		PlaidAccessToken: &plaidAccessToken,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
	}

	bankAccountCreated := models.BankAccount{
		BankAccountID: bankAccountID,
		UserID:   &userID,
		Name: &name,
		PlaidAccessToken: &plaidAccessToken,
		PlaidAccountID: &plaidAccountID,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	brokerageAccountID := uuid.New()
	alpacaAccountID := "alpacaAccountID"

	brokerageAccount := models.BrokerageAccount{
		BrokerageAccountID: brokerageAccountID,
		UserID:   &userID,
		AlpacaAccountID: &alpacaAccountID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonStr := []byte(fmt.Sprintf(
		`{
				"user_id": "%s",
				"brokerage_account_id": "%s",
				"plaid_public_token": "%s"
			}`,
		userID.String(),
		brokerageAccountID,
		plaidPublicToken,
	))

	testServer.Store.On("GetBrokerageAccount", brokerageAccountID).Return(&brokerageAccount, nil)
	testServer.PlaidClient.On("GetAccessToken", plaidPublicToken).Return(plaidAccessToken, nil)
	testServer.PlaidClient.On("GetAccount", plaidAccessToken).Return(plaidAccountID, name, nil)
	testServer.PlaidClient.On(
		"CreateProcessorToken",
		plaidAccessToken,
		plaidAccountID,
		"alpaca",
	).Return(plaidProcessorToken, nil)
	testServer.Broker.On(
		"CreateACHRelationship",
		alpacaAccountID,
		plaidProcessorToken,
	).Return(alpacaACHRelationshipID, nil)
	testServer.Store.On("CreateBankAccount", bankAccountPayload).Return(&bankAccountCreated, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/bank-accounts",
		bytes.NewBuffer(jsonStr),
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().CreateBankAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankAccountResponse models.BankAccount
	errDecoder := decoder.Decode(&bankAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccountCreated.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccountCreated.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccountCreated.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccountCreated.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccountCreated.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccountCreated.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccountCreated.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccountCreated.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
	testServer.PlaidClient.AssertExpectations(t)
	testServer.Broker.AssertExpectations(t)
}

func TestBankAccountModify(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankAccountID := uuid.New()
	userID := uuid.New()

	name := "name"
	plaidAccessToken := "plaidAccessToken"
	plaidAccountID := "plaidAccountID"
	alpacaACHRelationshipID := "alpacaACHRelationshipID"

	jsonStr := []byte(`{}`)

	bankAccountPayload := models.BankAccount{}

	bankAccountModified := models.BankAccount{
		BankAccountID: bankAccountID,
		UserID:   &userID,
		Name: &name,
		PlaidAccessToken: &plaidAccessToken,
		PlaidAccountID: &plaidAccountID,
		AlpacaACHRelationshipID: &alpacaACHRelationshipID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ModifyBankAccount", bankAccountModified.BankAccountID, bankAccountPayload).Return(&bankAccountModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/bank-accounts",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_account_id", bankAccountModified.BankAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ModifyBankAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankAccountResponse models.BankAccount
	errDecoder := decoder.Decode(&bankAccountResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankAccountResponse.BankAccountID, bankAccountModified.BankAccountID)
	assert.Equal(t, *bankAccountResponse.UserID, *bankAccountModified.UserID)
	assert.Equal(t, *bankAccountResponse.Name, *bankAccountModified.Name)
	assert.Equal(t, *bankAccountResponse.PlaidAccessToken, *bankAccountModified.PlaidAccessToken)
	assert.Equal(t, *bankAccountResponse.PlaidAccountID, *bankAccountModified.PlaidAccountID)
	assert.Equal(t, *bankAccountResponse.AlpacaACHRelationshipID, *bankAccountModified.AlpacaACHRelationshipID)
	assert.WithinDuration(t, bankAccountResponse.CreatedAt, bankAccountModified.CreatedAt, 0)
	assert.WithinDuration(t, bankAccountResponse.UpdatedAt, bankAccountModified.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankAccountDelete(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankAccountID := uuid.New()

	testServer.Store.On("DeleteBankAccount", bankAccountID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/bank-accounts",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_account_id", bankAccountID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().DeleteBankAccount(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	testServer.Store.AssertExpectations(t)
}
