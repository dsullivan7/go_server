package controllers_test

import (
	"time"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/models"
	"go_server/test/unit/internal/server/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestPortfolioGet(t *testing.T) {
	t.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(t, err)

	portfolioID := uuid.New()
	userID := uuid.New()

	portfolio := models.Portfolio{
		PortfolioID: portfolioID,
		UserID:   &userID,
		Risk:   3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStore.On("GetPortfolio", portfolioID).Return(&portfolio, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/portfolios",
		nil,
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portfolio_id", portfolioID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	controllers.GetPortfolio(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var portfolioResponse models.Portfolio
	errDecoder := decoder.Decode(&portfolioResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, portfolioResponse.PortfolioID, portfolio.PortfolioID)
	assert.Equal(t, portfolioResponse.UserID, portfolio.UserID)
	assert.Equal(t, portfolioResponse.Risk, portfolio.Risk)
	assert.WithinDuration(t, portfolioResponse.CreatedAt, portfolio.CreatedAt, 0)
	assert.WithinDuration(t, portfolioResponse.UpdatedAt, portfolio.UpdatedAt, 0)

	mockStore.AssertExpectations(t)
}

func TestPortfolioList(t *testing.T) {
	t.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(t, err)

	portfolioID1 := uuid.New()
	userID1 := uuid.New()

	portfolio1 := models.Portfolio{
		PortfolioID: portfolioID1,
		UserID:   &userID1,
		Risk:   3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	portfolioID2 := uuid.New()
	userID2 := uuid.New()

	portfolio2 := models.Portfolio{
		PortfolioID: portfolioID2,
		UserID:   &userID2,
		Risk:   4,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStore.On("ListPortfolios", map[string]interface{}{}).Return([]models.Portfolio{portfolio1, portfolio2}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/users",
		nil,
	)

	w := httptest.NewRecorder()

	controllers.ListPortfolios(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var portfoliosFound []models.Portfolio
	errDecoder := decoder.Decode(&portfoliosFound)
	assert.Nil(t, errDecoder)

	assert.Equal(t, 2, len(portfoliosFound))

	var portfolioResponse models.Portfolio

	for _, value := range portfoliosFound {
		if value.PortfolioID == portfolio1.PortfolioID {
			portfolioResponse = value

			break
		}
	}

	assert.Equal(t, portfolioResponse.PortfolioID, portfolio1.PortfolioID)
	assert.Equal(t, portfolioResponse.UserID, portfolio1.UserID)
	assert.Equal(t, portfolioResponse.Risk, portfolio1.Risk)
	assert.WithinDuration(t, portfolioResponse.CreatedAt, portfolio1.CreatedAt, 0)
	assert.WithinDuration(t, portfolioResponse.UpdatedAt, portfolio1.UpdatedAt, 0)

	for _, value := range portfoliosFound {
		if value.PortfolioID == portfolio2.PortfolioID {
			portfolioResponse = value

			break
		}
	}

	assert.Equal(t, portfolioResponse.PortfolioID, portfolio2.PortfolioID)
	assert.Equal(t, portfolioResponse.UserID, portfolio2.UserID)
	assert.Equal(t, portfolioResponse.Risk, portfolio2.Risk)
	assert.WithinDuration(t, portfolioResponse.CreatedAt, portfolio2.CreatedAt, 0)
	assert.WithinDuration(t, portfolioResponse.UpdatedAt, portfolio2.UpdatedAt, 0)

	mockStore.AssertExpectations(t)
}

func TestPortfolioCreate(t *testing.T) {
	t.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(t, err)

	userID := uuid.New()
	risk := 3

	jsonStr := []byte(fmt.Sprintf(
		`{
				"user_id": "%s",
				"risk": %d
			}`,
		userID.String(),
		risk,
	))

	portfolioPayload := models.Portfolio{
		UserID:   &userID,
		Risk:  risk,
	}

	portfolioCreated := models.Portfolio{
		PortfolioID:    uuid.New(),
		UserID:   &userID,
		Risk:  risk,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStore.On("CreatePortfolio", portfolioPayload).Return(&portfolioCreated, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/portfolios",
		bytes.NewBuffer(jsonStr),
	)

	w := httptest.NewRecorder()

	controllers.CreatePortfolio(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var portfolioResponse models.Portfolio
	errDecoder := decoder.Decode(&portfolioResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, portfolioResponse.PortfolioID, portfolioCreated.PortfolioID)
	assert.Equal(t, portfolioResponse.UserID, portfolioCreated.UserID)
	assert.Equal(t, portfolioResponse.Risk, portfolioCreated.Risk)
	assert.WithinDuration(t, portfolioResponse.CreatedAt, portfolioCreated.CreatedAt, 0)
	assert.WithinDuration(t, portfolioResponse.UpdatedAt, portfolioCreated.UpdatedAt, 0)

	mockStore.AssertExpectations(t)
}

func TestPortfolioModify(t *testing.T) {
	t.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(t, err)

	userID := uuid.New()
	risk := 3

	jsonStr := []byte(fmt.Sprintf(
		`{
				"risk": %d
			}`,
		risk,
	))

	portfolioPayload := models.Portfolio{
		Risk:  risk,
	}

	portfolioModified := models.Portfolio{
		PortfolioID:    uuid.New(),
		UserID:   &userID,
		Risk:  risk,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStore.On("ModifyPortfolio", portfolioModified.PortfolioID, portfolioPayload).Return(&portfolioModified, nil)

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/portfolios",
		bytes.NewBuffer(jsonStr),
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portfolio_id", portfolioModified.PortfolioID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	controllers.ModifyPortfolio(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))

	decoder := json.NewDecoder(res.Body)

	var portfolioResponse models.Portfolio
	errDecoder := decoder.Decode(&portfolioResponse)
	assert.Nil(t, errDecoder)

	assert.Equal(t, portfolioResponse.PortfolioID, portfolioModified.PortfolioID)
	assert.Equal(t, portfolioResponse.UserID, portfolioModified.UserID)
	assert.Equal(t, portfolioResponse.Risk, portfolioModified.Risk)
	assert.WithinDuration(t, portfolioResponse.CreatedAt, portfolioModified.CreatedAt, 0)
	assert.WithinDuration(t, portfolioResponse.UpdatedAt, portfolioModified.UpdatedAt, 0)

	mockStore.AssertExpectations(t)
}

func TestPortfolioDelete(t *testing.T) {
	t.Parallel()

	controllers, mockStore, err := controllers.Setup()
	assert.Nil(t, err)

	portfolioID := uuid.New()

	mockStore.On("DeletePortfolio", portfolioID).Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/portfolios",
		nil,
	)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("portfolio_id", portfolioID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	controllers.DeletePortfolio(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	mockStore.AssertExpectations(t)
}
