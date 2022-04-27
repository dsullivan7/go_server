package controllers

import (
	"fmt"
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"crypto/rand"
	"encoding/hex"
)

const idLength = 16
const secretLength = 32

func (c *Controllers) getGroupResponse(
	group models.Group,
) (*models.Group, error) {
	groupResponse := models.Group(group)

	apiClientKey, errDecryptKey := c.cipher.Decrypt(group.APIClientKey, c.config.EncryptionKey)

	if errDecryptKey != nil {
		return nil, fmt.Errorf("error decrypting client_id: %w", errDecryptKey)
	}

	apiClientSecret, errDecryptSecret := c.cipher.Decrypt(group.APIClientSecret, c.config.EncryptionKey)

	if errDecryptSecret != nil {
		return nil, fmt.Errorf("error decrypting client_secret: %w", errDecryptSecret)
	}

	groupResponse.APIClientKey = apiClientKey
	groupResponse.APIClientSecret = apiClientSecret

	return &groupResponse, nil
}

func randomHex(n int) (string, error) {
  bytes := make([]byte, n)
  if _, err := rand.Read(bytes); err != nil {
    return "", err
  }

  return hex.EncodeToString(bytes), nil
}

func (c *Controllers) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := uuid.Must(uuid.Parse(chi.URLParam(r, "group_id")))

	group, err := c.store.GetGroup(groupID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	groupResponse, errResponse := c.getGroupResponse(*group)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.JSON(w, r, groupResponse)
}

func (c *Controllers) ListGroups(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}

	groups, err := c.store.ListGroups(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, groups)
}

func (c *Controllers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupRequest map[string]interface{}

	errDecode := json.NewDecoder(r.Body).Decode(&groupRequest)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	apiClientKey, errHexID := randomHex(idLength)

	if errHexID != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errHexID})

		return
	}

	apiClientSecret, errHexSecret := randomHex(secretLength)

	if errHexSecret != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errHexSecret})

		return
	}

	apiClientKeyEnc, errEncryptKey := c.cipher.Encrypt(apiClientKey, c.config.EncryptionKey)

	if errEncryptKey != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEncryptKey})

		return
	}

	apiClientSecretEnc, errEncryptSecret := c.cipher.Encrypt(apiClientSecret, c.config.EncryptionKey)

	if errEncryptSecret != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEncryptSecret})

		return
	}

	groupPayload := models.Group{
		Name: groupRequest["name"].(string),
		APIClientKey: apiClientKeyEnc,
		APIClientSecret: apiClientSecretEnc,
	}

	group, err := c.store.CreateGroup(groupPayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	if (groupRequest["group_users"] != nil) {
		for _, groupUserPayload := range groupRequest["group_users"].([]interface{}) {
			userID := uuid.Must(uuid.Parse(groupUserPayload.(map[string]interface{})["user_id"].(string)))
			groupUser := models.GroupUser{
				UserID: userID,
				GroupID: group.GroupID,
			}

			c.store.CreateGroupUser(groupUser)
		}
	}

	groupResponse, errResponse := c.getGroupResponse(*group)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, groupResponse)
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

	groupResponse, errResponse := c.getGroupResponse(*group)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.JSON(w, r, groupResponse)
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
