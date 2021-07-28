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
	reviews := services.ListReviews(&models.Review{})

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
