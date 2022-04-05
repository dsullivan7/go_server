package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go_server/internal/errors"
	"go_server/internal/models"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

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

	if errLogin != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errLogin})

		return
	}

	defer resLogin.Body.Close()

	urlPayments := "https://a069-access.nyc.gov/Rest/v1/ua/anyc/payments/1"

	reqPayments, errReqPayments := http.NewRequestWithContext(
		context,
		http.MethodGet,
		urlPayments,
		nil,
	)

	if errReqPayments != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReqPayments})

		return
	}

	reqPayments.Header.Add("Referer", "https://a069-access.nyc.gov/accesshra/anycuserhome")

	resPayments, errResPayments := client.Do(reqPayments)

	if errResPayments != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errResPayments})

		return
	}

	defer resPayments.Body.Close()
	paymentsB, errReadPayments := io.ReadAll(resPayments.Body)

	if errReadPayments != nil {
		c.utils.HandleError(w, r, errors.HTTPUserError{Err: errReadPayments})

		return
	}

	res := map[string]string{
		"response": string(paymentsB),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, res)
}
