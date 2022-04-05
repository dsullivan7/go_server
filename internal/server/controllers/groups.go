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

func (c *Controllers) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_id")))

	group, err := c.store.GetGroup(groupID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, group)
}

func (c *Controllers) ListGroups(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	groups, err := c.store.ListGroups(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, groups)
}

func (c *Controllers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupPayload models.Group

	errDecode := json.NewDecoder(r.Body).Decode(&groupPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	group, err := c.store.CreateGroup(groupPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, group)
}

func (c *Controllers) ModifyGroup(w http.ResponseWriter, r *http.Request) {
	var groupPayload models.Group

	groupID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&groupPayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	group, err := c.store.ModifyGroup(groupID, groupPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, group)
}

func (c *Controllers) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_id")))

	err := c.store.DeleteGroup(groupID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
