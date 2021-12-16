package twocaptcha

import (
	"context"
	"fmt"
	"go_server/internal/captcha"
	"go_server/internal/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const TwoCaptchaRequestToCompleteURL = "http://2captcha.com/in.php?key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s"
const TwoCaptchaResponseURL = "http://2captcha.com/res.php?key=%s&action=get&id=%s"
const TwoCaptchaNotReadyResponse = "CAPCHA_NOT_READY"
const TwoCaptchaFirstWait = 15
const TwoCaptchaRetry = 5

type Captcha struct {
	captchaKey string
	logger     logger.Logger
}

func NewTwoCaptcha(captchaKey string, logger logger.Logger) captcha.Captcha {
	return &Captcha{
		captchaKey: captchaKey,
		logger:     logger,
	}
}

func (captcha *Captcha) createReCaptchaV2Request(googleKey string, url string) (*string, error) {
	context := context.Background()
	captchaURL := fmt.Sprintf(TwoCaptchaRequestToCompleteURL, captcha.captchaKey, googleKey, url)

	// send the request to 2captcha to complete a captcha
	req, errReq := http.NewRequestWithContext(
		context,
		http.MethodGet,
		captchaURL,
		nil,
	)

	if errReq != nil {
		return nil, fmt.Errorf("failed to send request to complete captcha: %w", errReq)
	}

	res, errRes := http.DefaultClient.Do(req)

	if errRes != nil {
		return nil, fmt.Errorf("failed to get response to complete captcha: %w", errRes)
	}

	defer res.Body.Close()

	captcha.logger.Info("sent request to 2captcha")

	bodyToComplete, errReadToComplete := ioutil.ReadAll(res.Body)

	if errReadToComplete != nil {
		return nil, fmt.Errorf("failed to read response to complete captcha: %w", errReadToComplete)
	}

	// extract the id to use later to check the status of the request
	captchaArray := strings.Split(string(bodyToComplete), "|")
	captchaID := captchaArray[1]

	return &captchaID, nil
}

func (captcha *Captcha) getReCaptchaV2Response(captchaID string) (*string, error) {
	context := context.Background()
	captchaCompleteURL := fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=get&id=%s", captcha.captchaKey, captchaID)

	// timeout for 20 seconds to wait for the response
	time.Sleep(TwoCaptchaFirstWait * time.Second)

	var captchaComplete *string

	for captchaComplete == nil {
		// ping 2captcha every 5 seconds to determine if the request has finished
		time.Sleep(TwoCaptchaRetry * time.Second)

		req, errReq := http.NewRequestWithContext(
			context,
			http.MethodGet,
			captchaCompleteURL,
			nil,
		)

		if errReq != nil {
			return nil, fmt.Errorf("failed to send request to check captcha: %w", errReq)
		}

		res, errRes := http.DefaultClient.Do(req)

		if errRes != nil {
			return nil, fmt.Errorf("failed to get response to check captcha: %w", errRes)
		}

		defer res.Body.Close()

		bodyComplete, errReadComplete := ioutil.ReadAll(res.Body)

		if errReadComplete != nil {
			return nil, fmt.Errorf("failed to read captcha response: %w", errReadComplete)
		}

		bodyCompleteString := string(bodyComplete)

		if bodyCompleteString != TwoCaptchaNotReadyResponse {
			captchaCompleteArray := strings.Split(bodyCompleteString, "|")
			captchaComplete = &captchaCompleteArray[1]
		}
	}

	return captchaComplete, nil
}

func (captcha *Captcha) SolveReCaptchaV2(googleKey string, url string) (*string, error) {
	captcha.logger.Info("solving recaptcha v2")

	captchaID, errCaptchaID := captcha.createReCaptchaV2Request(googleKey, url)
	if errCaptchaID != nil {
		return nil, errCaptchaID
	}

	captchaResponse, errCaptchaRespose := captcha.getReCaptchaV2Response(*captchaID)
	if errCaptchaRespose != nil {
		return nil, errCaptchaRespose
	}

	return captchaResponse, nil
}
