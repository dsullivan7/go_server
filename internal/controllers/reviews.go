package controllers

import (
	"encoding/json"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func (c *Controllers) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	review := c.store.GetReview(reviewID)

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

	reviews := c.store.ListReviews(query)

	render.JSON(w, r, reviews)
}

func (c *Controllers) CreateReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review

	err := json.NewDecoder(r.Body).Decode(&reviewPayload)
	if err != nil {
		w.WriteHeader(HTTP400)

		return
	}

	review := c.store.CreateReview(reviewPayload)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, review)
}

func (c *Controllers) ModifyReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review

	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	err := json.NewDecoder(r.Body).Decode(&reviewPayload)
	if err != nil {
		w.WriteHeader(HTTP400)

		return
	}

	review := c.store.ModifyReview(reviewID, reviewPayload)

	render.JSON(w, r, review)
}

func (c *Controllers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	c.store.DeleteReview(reviewID)

	w.WriteHeader(http.StatusNoContent)
}
