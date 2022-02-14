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

func TestBankTransferGet(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankTransferID := uuid.New()
	userID := uuid.New()

	bankTransfer := models.BankTransfer{
		BankTransferID: bankTransferID,
		UserID:   &userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("GetBankTransfer", bankTransferID).Return(&bankTransfer, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/bank-transfers",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_transfer_id", bankTransferID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().GetBankTransfer(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankTransferResponse models.BankTransfer
	errDecoder := decoder.Decode(&bankTransferResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransfer.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransfer.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransfer.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransfer.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankTransferList(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankTransferID1 := uuid.New()
	userID1 := uuid.New()

	bankTransfer1 := models.BankTransfer{
		BankTransferID: bankTransferID1,
		UserID:   &userID1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	bankTransferID2 := uuid.New()
	userID2 := uuid.New()

	bankTransfer2 := models.BankTransfer{
		BankTransferID: bankTransferID2,
		UserID:   &userID2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ListBankTransfers", map[string]interface{}{}).Return([]models.BankTransfer{bankTransfer1, bankTransfer2}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/bank-transfers",
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBankTransfers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankTransfersFound []models.BankTransfer
	errDecoder := decoder.Decode(&bankTransfersFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(bankTransfersFound))

	var bankTransferResponse models.BankTransfer

	for _, value := range bankTransfersFound {
		if value.BankTransferID == bankTransfer1.BankTransferID {
			bankTransferResponse = value

			break
		}
	}

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransfer1.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransfer1.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransfer1.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransfer1.UpdatedAt, 0)

	for _, value := range bankTransfersFound {
		if value.BankTransferID == bankTransfer2.BankTransferID {
			bankTransferResponse = value

			break
		}
	}

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransfer2.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransfer2.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransfer2.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransfer2.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankTransferListQueryParams(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankTransferID := uuid.New()
	userID := uuid.New()

	bankTransfer := models.BankTransfer{
		BankTransferID: bankTransferID,
		UserID:   &userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ListBankTransfers", map[string]interface{}{ "user_id": userID.String() }).Return([]models.BankTransfer{bankTransfer}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprint("/api/bank-transfers?user_id=", userID),
		nil,
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ListBankTransfers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankTransfersFound []models.BankTransfer
	errDecoder := decoder.Decode(&bankTransfersFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 1, len(bankTransfersFound))

	bankTransferResponse := bankTransfersFound[0]

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransfer.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransfer.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransfer.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransfer.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankTransferCreate(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()

	jsonStr := []byte(fmt.Sprintf(
		`{
				"user_id": "%s"
			}`,
		userID.String(),
	))

	bankTransferPayload := models.BankTransfer{
		UserID:   &userID,
	}

	bankTransferCreated := models.BankTransfer{
		BankTransferID:    uuid.New(),
		UserID:   &userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("CreateBankTransfer", bankTransferPayload).Return(&bankTransferCreated, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/bank-transfers",
		bytes.NewBuffer(jsonStr),
	)

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().CreateBankTransfer(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankTransferResponse models.BankTransfer
	errDecoder := decoder.Decode(&bankTransferResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransferCreated.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransferCreated.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransferCreated.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransferCreated.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankTransferModify(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	userID := uuid.New()

	jsonStr := []byte(`{}`)

	bankTransferPayload := models.BankTransfer{}

	bankTransferModified := models.BankTransfer{
		BankTransferID:    uuid.New(),
		UserID:   &userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testServer.Store.On("ModifyBankTransfer", bankTransferModified.BankTransferID, bankTransferPayload).Return(&bankTransferModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/bank-transfers",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_transfer_id", bankTransferModified.BankTransferID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().ModifyBankTransfer(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var bankTransferResponse models.BankTransfer
	errDecoder := decoder.Decode(&bankTransferResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, bankTransferResponse.BankTransferID, bankTransferModified.BankTransferID)
	assert.Equal(t, bankTransferResponse.UserID, bankTransferModified.UserID)
	assert.WithinDuration(t, bankTransferResponse.CreatedAt, bankTransferModified.CreatedAt, 0)
	assert.WithinDuration(t, bankTransferResponse.UpdatedAt, bankTransferModified.UpdatedAt, 0)

	testServer.Store.AssertExpectations(t)
}

func TestBankTransferDelete(t *testing.T) {
	t.Parallel()

	testServer, err := utils.NewTestServer()
	assert.Nil(t, err)

	bankTransferID := uuid.New()

	testServer.Store.On("DeleteBankTransfer", bankTransferID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/bank-transfers",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bank_transfer_id", bankTransferID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	testServer.Server.GetControllers().DeleteBankTransfer(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	testServer.Store.AssertExpectations(t)
}
