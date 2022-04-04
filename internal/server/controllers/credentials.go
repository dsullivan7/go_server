package controllers

import (
	"fmt"
	"encoding/json"
	"go_server/internal/errors"
	"go_server/internal/models"
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"io"

	"github.com/go-chi/render"
)

//nolint:funlen
func (c *Controllers) CreateCredential(w http.ResponseWriter, r *http.Request) {
	var credentialPayload models.Credential

	errDecode := json.NewDecoder(r.Body).Decode(&credentialPayload)
	if errDecode != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errDecode})

		return
	}

	context := context.Background()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
    Jar: jar,
	}

	urlToLogin := "https://a069-access.nyc.gov/Rest/j_security_check"

	data := url.Values{}
	data.Set("j_username", credentialPayload.Username)
	data.Set("j_password", credentialPayload.Password)
	data.Set("user_type", fmt.Sprint("EXTERNAL;", r.RemoteAddr))

	reqLogin, _ := http.NewRequestWithContext(
		context,
		http.MethodPost,
		urlToLogin,
		strings.NewReader(data.Encode()),
	)

	reqLogin.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resLogin, errLogin := client.Do(reqLogin)

	if (errLogin != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errLogin})

		return
	}

	defer resLogin.Body.Close()

	urlProfile := "https://a069-access.nyc.gov/Rest/v1/ua/profile"

	reqProfile, errReqProfile := http.NewRequestWithContext(
		context,
		http.MethodGet,
		urlProfile,
		nil,
	)

	if (errReqProfile != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReqProfile})

		return
	}

	reqProfile.Header.Add("Referer", "https://a069-access.nyc.gov/accesshra/anycuserhome")

	resProfile, errResPorfile := client.Do(reqProfile)

	if (errResPorfile != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errResPorfile})

		return
	}

	defer resProfile.Body.Close()
	profileB, errReadProfile := io.ReadAll(resProfile.Body)

	if (errReadProfile != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReadProfile})

		return
	}

	urlBenefits := "https://a069-access.nyc.gov/Rest/v1/ua/anyc/benefits"

	reqBenefits, errReqBenefits := http.NewRequestWithContext(
		context,
		http.MethodGet,
		urlBenefits,
		nil,
	)

	if (errReqBenefits != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReqBenefits})

		return
	}

	reqBenefits.Header.Add("Referer", "https://a069-access.nyc.gov/accesshra/anycuserhome")

	resBenefits, errResBenefits := client.Do(reqBenefits)
	if (errResBenefits != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errResBenefits})

		return
	}

	defer resBenefits.Body.Close()
	benefitsB, errReadBenefits := io.ReadAll(resBenefits.Body)

	if (errReadBenefits != nil) {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReadBenefits})

		return
	}

	res := map[string]string{
		"response": fmt.Sprint(string(profileB), "\n", string(benefitsB)),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}
