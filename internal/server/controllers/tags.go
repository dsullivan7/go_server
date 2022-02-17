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

func (c *Controllers) GetTag(w http.ResponseWriter, r *http.Request) {
	tagID := uuid.Must(uuid.Parse(chi.URLParam(r, "tag_id")))

	tag, err := c.store.GetTag(tagID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, tag)
}

func (c *Controllers) ListTags(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	tags, err := c.store.ListTags(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, tags)
}

func (c *Controllers) CreateTag(w http.ResponseWriter, r *http.Request) {
	var tagPayload models.Tag

	errDecode := json.NewDecoder(r.Body).Decode(&tagPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	tag, err := c.store.CreateTag(tagPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, tag)
}

func (c *Controllers) ModifyTag(w http.ResponseWriter, r *http.Request) {
	var tagPayload models.Tag

	tagID := uuid.Must(uuid.Parse(chi.URLParam(r, "tag_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&tagPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	tag, err := c.store.ModifyTag(tagID, tagPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, tag)
}

func (c *Controllers) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagID := uuid.Must(uuid.Parse(chi.URLParam(r, "tag_id")))

	err := c.store.DeleteTag(tagID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
