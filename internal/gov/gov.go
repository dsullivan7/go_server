package gov

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Profile struct {
	EBTSNAPBalance string
}

type IGov interface {
	GetProfile(username string, password string, ipAddress string, portalType string) (*Profile, error)
}

type Gov struct{}

func NewGov() IGov {
	return &Gov{}
}

func (gv *Gov) GetProfile(username string, password string, ipAddress string, portalType string) (*Profile, error) {
	if portalType == "accesshra" {
		return gv.getAccessHRAProfile(username, password, ipAddress)
	}

	return nil, fmt.Errorf("Portal type not found")
}

// nolint:funlen
func (gv *Gov) getAccessHRAProfile(username string, password string, ipAddress string) (*Profile, error) {
	context := context.Background()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	urlToLogin := "https://a069-access.nyc.gov/Rest/j_security_check"

	data := url.Values{}
	data.Set("j_username", username)
	data.Set("j_password", password)
	data.Set("user_type", fmt.Sprint("EXTERNAL;", ipAddress))

	reqLogin, _ := http.NewRequestWithContext(
		context,
		http.MethodPost,
		urlToLogin,
		strings.NewReader(data.Encode()),
	)

	reqLogin.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resLogin, errLogin := client.Do(reqLogin)

	if errLogin != nil {
		return nil, fmt.Errorf("failed to login to accessHRA: %w", errLogin)
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
		return nil, fmt.Errorf("failed to create payments request: %w", errReqPayments)
	}

	reqPayments.Header.Add("Referer", "https://a069-access.nyc.gov/accesshra/anycuserhome")

	resPayments, errResPayments := client.Do(reqPayments)

	if errResPayments != nil {
		return nil, fmt.Errorf("failed to get payments response: %w", errResPayments)
	}

	defer resPayments.Body.Close()

	decoder := json.NewDecoder(resPayments.Body)

	var accessHRAResponse interface{}

	errDecode := decoder.Decode(&accessHRAResponse)

	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", errDecode)
	}

	return &Profile{EBTSNAPBalance: accessHRAResponse.(map[string]interface{})["snapEBTBalance"].(string)}, nil
}
