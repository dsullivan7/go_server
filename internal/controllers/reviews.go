package controllers

import (
	"net/http"
	"encoding/json"

	"go_server/internal/services"
	"go_server/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	review := services.GetReview(reviewID)

	render.JSON(w, r, review)
}

func ListReviews(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	fromUserID := r.URL.Query().Get("from_user_id")
	toUserID := r.URL.Query().Get("to_user_id")

	if (fromUserID != "") {
		query["from_user_id"] = fromUserID
	}

	if (toUserID != "") {
		query["to_user_id"] = toUserID
	}

	reviews := services.ListReviews(query)

	render.JSON(w, r, reviews)
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review

	json.NewDecoder(r.Body).Decode(&reviewPayload)

	review := services.CreateReview(reviewPayload)

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, review)
}

func ModifyReview(w http.ResponseWriter, r *http.Request) {
	var reviewPayload models.Review
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	json.NewDecoder(r.Body).Decode(&reviewPayload)

	review := services.ModifyReview(reviewID, reviewPayload)

	render.JSON(w, r, review)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := uuid.Must(uuid.Parse(chi.URLParam(r, "reviewID")))

	services.DeleteReview(reviewID)

	w.WriteHeader(http.StatusNoContent)
}
