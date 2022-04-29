package controllers

import (
	"encoding/json"
	"fmt"
	"go_server/internal/errors"
	"go_server/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ProfileResponse struct {
	ProfileID      uuid.UUID `json:"profile_id"`
	EBTSNAPBalance string    `json:"ebt_snap_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (c *Controllers) getProfileResponse(
	profile models.Profile,
	ipAddress string,
) (*ProfileResponse, error) {
	var profileResponse ProfileResponse

	profileResponse.ProfileID = profile.ProfileID
	profileResponse.CreatedAt = profile.CreatedAt
	profileResponse.UpdatedAt = profile.UpdatedAt

	username, errDecryptUser := c.cipher.Decrypt(profile.Username, c.config.EncryptionKey)

	if errDecryptUser != nil {
		return nil, fmt.Errorf("error decrypting username: %w", errDecryptUser)
	}

	password, errDecryptPass := c.cipher.Decrypt(profile.Password, c.config.EncryptionKey)

	if errDecryptPass != nil {
		return nil, fmt.Errorf("error decrypting password: %w", errDecryptPass)
	}

	benefitsProfile, errBenefits := c.gov.GetProfile(username, password, ipAddress, profile.Type)

	if errBenefits != nil {
		return nil, fmt.Errorf("error retrieving benefits: %w", errBenefits)
	}

	profileResponse.EBTSNAPBalance = benefitsProfile.EBTSNAPBalance

	return &profileResponse, nil
}

func (c *Controllers) GetProfile(w http.ResponseWriter, r *http.Request) {
	profileID := uuid.Must(uuid.Parse(chi.URLParam(r, "profile_id")))

	profile, err := c.store.GetProfile(profileID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	profileResponse, errResponse := c.getProfileResponse(*profile, r.RemoteAddr)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, profileResponse)
}

func (c *Controllers) ListProfiles(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{}
	profileID := r.URL.Query().Get("profile_id")

	if profileID != "" {
		query["profile_id"] = profileID
	}

	profiles, err := c.store.ListProfiles(query)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.JSON(w, r, profiles)
}

func (c *Controllers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profilePayload models.Profile

	errDecode := json.NewDecoder(r.Body).Decode(&profilePayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	usernameEncrypted, errEncryptUser := c.cipher.Encrypt(profilePayload.Username, c.config.EncryptionKey)

	if errEncryptUser != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEncryptUser})

		return
	}

	passwordEncrypted, errEncryptPass := c.cipher.Encrypt(profilePayload.Password, c.config.EncryptionKey)

	if errEncryptPass != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errEncryptPass})

		return
	}

	profilePayload.Username = usernameEncrypted
	profilePayload.Password = passwordEncrypted

	profile, err := c.store.CreateProfile(profilePayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	profileResponse, errResponse := c.getProfileResponse(*profile, r.RemoteAddr)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, profileResponse)
}

func (c *Controllers) ModifyProfile(w http.ResponseWriter, r *http.Request) {
	var profilePayload models.Profile

	profileID := uuid.Must(uuid.Parse(chi.URLParam(r, "profile_id")))

	errDecode := json.NewDecoder(r.Body).Decode(&profilePayload)

	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	profile, err := c.store.ModifyProfile(profileID, profilePayload)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	profileResponse, errResponse := c.getProfileResponse(*profile, r.RemoteAddr)

	if errResponse != nil {
		c.utils.HandleError(w, r, errors.HTTPServerError{Err: errResponse})

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, profileResponse)
}

func (c *Controllers) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileID := uuid.Must(uuid.Parse(chi.URLParam(r, "profile_id")))

	err := c.store.DeleteProfile(profileID)

	if err != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: err})

		return
	}

	render.NoContent(w, r)
}
