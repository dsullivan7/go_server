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

func (c *Controllers) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	review, err := c.store.GetReview(reviewID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, review)
}

func (c *Controllers) ListReviews(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	fromUserID := r.URL.Query().Get("from_user_id")
	toUserID := r.URL.Query().Get("to_user_id")

	if fromUserID != "" {
		query["from_user_id"] = fromUserID
	}

	if toUserID != "" {
		query["to_user_id"] = toUserID
	}

	reviews, err := c.store.ListReviews(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, reviews)
}

func (c *Controllers) CreateReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review

	errDecode := json.NewDecoder(r.Body).Decode(&reviewPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	review, err := c.store.CreateReview(reviewPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, review)
}

func (c *Controllers) ModifyReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review

	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	errDecode := json.NewDecoder(r.Body).Decode(&reviewPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	review, err := c.store.ModifyReview(reviewID, reviewPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, review)
}

func (c *Controllers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	err := c.store.DeleteReview(reviewID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
