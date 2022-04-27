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

func (c *Controllers) GetGroupUser(w http.ResponseWriter, r *http.Request) {
	groupUserID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_user_id")))

	groupUser, err := c.store.GetGroupUser(groupUserID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, groupUser)
}

func (c *Controllers) ListGroupUsers(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	userID := r.URL.Query().Get("user_id")

	if userID != "" {
		query["user_id"] = userID
	}

	groupUsers, err := c.store.ListGroupUsers(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, groupUsers)
}

func (c *Controllers) CreateGroupUser(w http.ResponseWriter, r *http.Request) {
	var groupUserPayload models.GroupUser

	errDecode := json.NewDecoder(r.Body).Decode(&groupUserPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	groupUser, err := c.store.CreateGroupUser(groupUserPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, groupUser)
}

func (c *Controllers) ModifyGroupUser(w http.ResponseWriter, r *http.Request) {
	var groupUserPayload models.GroupUser

	groupUserID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_user_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&groupUserPayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	groupUser, err := c.store.ModifyGroupUser(groupUserID, groupUserPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, groupUser)
}

func (c *Controllers) DeleteGroupUser(w http.ResponseWriter, r *http.Request) {
	groupUserID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_user_id")))

	err := c.store.DeleteGroupUser(groupUserID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
