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

func (c *Controllers) GetUserIndustry(w http.ResponseWriter, r *http.Request) {
	userIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "user_industry_id")))

	userIndustry, err := c.store.GetUserIndustry(userIndustryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, userIndustry)
}

func (c *Controllers) ListUserIndustries(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := c.utils.GetQueryParamUUID(r, "user_id")

	if userID != uuid.Nil {
		query["user_id"] = userID
	}

	userIndustries, err := c.store.ListUserIndustries(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, userIndustries)
}

func (c *Controllers) CreateUserIndustry(w http.ResponseWriter, r *http.Request) {
	var userIndustryPayload models.UserIndustry

	errDecode := json.NewDecoder(r.Body).Decode(&userIndustryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	userIndustry, err := c.store.CreateUserIndustry(userIndustryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, userIndustry)
}

func (c *Controllers) ModifyUserIndustry(w http.ResponseWriter, r *http.Request) {
	var userIndustryPayload models.UserIndustry

	userIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "user_industry_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&userIndustryPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	userIndustry, err := c.store.ModifyUserIndustry(userIndustryID, userIndustryPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, userIndustry)
}

func (c *Controllers) DeleteUserIndustry(w http.ResponseWriter, r *http.Request) {
	userIndustryID := uuid.Must(uuid.Parse(chi.URLParam(r, "user_industry_id")))

	err := c.store.DeleteUserIndustry(userIndustryID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
