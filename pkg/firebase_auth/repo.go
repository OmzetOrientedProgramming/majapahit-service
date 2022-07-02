package firebaseauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Repo contains all the function that available of this repo package
type Repo interface {
	SendOTP(params SendOTPParams) (*SendOTPResult, error)
	VerifyOTP(params VerifyOTPParams) (*VerifyOTPResult, error)
	GetUserDataFromToken(token string) (*UserDataFromToken, error)
}

type repo struct {
	identityToolkitURL string
	secureTokenURL     string
	APIKey             string
}

func (r repo) VerifyOTP(params VerifyOTPParams) (*VerifyOTPResult, error) {
	URL := fmt.Sprintf("%s/v1/accounts:signInWithPhoneNumber?key=%s", r.identityToolkitURL, r.APIKey)

	reqBody := map[string]interface{}{
		"sessionInfo": params.SessionInfo,
		"code":        params.Code,
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	resp, err := http.Post(URL, "", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	respBuffer, _ := io.ReadAll(resp.Body)

	var respSuccess VerifyOTPResult
	err = json.Unmarshal(respBuffer, &respSuccess)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	if respSuccess.IDToken == "" {
		var errorResponse ErrorFromFirebase

		err = json.Unmarshal(respBuffer, &errorResponse)
		if err != nil {
			return nil, errors.Wrap(ErrInternalServer, err.Error())
		}

		if strings.Contains(errorResponse.Error.Message, "SESSION_EXPIRED") {
			return nil, errors.Wrap(ErrInputValidation, "session expired")
		}

		if strings.Contains(errorResponse.Error.Message, "INVALID_CODE") {
			return nil, errors.Wrap(ErrInputValidation, "invalid OTP code")
		}

		if strings.Contains(errorResponse.Error.Message, "INVALID_SESSION_INFO") {
			return nil, errors.Wrap(ErrInputValidation, "invalid session")
		}

		return nil, errors.Wrap(ErrInternalServer, errorResponse.Error.Message)
	}

	return &respSuccess, nil

}

func (r repo) SendOTP(params SendOTPParams) (*SendOTPResult, error) {
	URL := fmt.Sprintf("%s/v1/accounts:sendVerificationCode?key=%s", r.identityToolkitURL, r.APIKey)

	reqBody := map[string]interface{}{
		"phoneNumber":    params.PhoneNumber,
		"recaptchaToken": params.RecaptchaToken,
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	resp, err := http.Post(URL, "", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	respBuffer, _ := io.ReadAll(resp.Body)

	var respAbstract map[string]interface{}
	_ = json.Unmarshal(respBuffer, &respAbstract)

	var (
		ok                bool
		sessionInfoString string
		errorResponse     ErrorFromFirebase
	)

	if _, ok = respAbstract["sessionInfo"]; !ok {
		err = json.Unmarshal(respBuffer, &errorResponse)
		if err != nil {
			return nil, errors.Wrap(ErrInternalServer, err.Error())
		}

		if strings.Contains(errorResponse.Error.Message, "CAPTCHA_CHECK_FAILED : Recaptcha verification failed - EXPIRED") {
			return nil, errors.Wrap(ErrInputValidation, "recaptcha expired")
		}

		if strings.Contains(errorResponse.Error.Message, "TOO_MANY_ATTEMPTS_TRY_LATER") {
			return nil, errors.Wrap(ErrInputValidation, "too many attempts. try again later")
		}

		return nil, errors.Wrap(ErrInternalServer, errorResponse.Error.Message)
	}

	if sessionInfoString, ok = respAbstract["sessionInfo"].(string); !ok {
		return nil, errors.Wrap(ErrInternalServer, "invalid session info format")
	}

	if sessionInfoString == "" {
		return nil, errors.Wrap(ErrInternalServer, "session info is empty")
	}

	return &SendOTPResult{
		SessionInfo: sessionInfoString,
	}, nil
}

func (r repo) GetUserDataFromToken(token string) (*UserDataFromToken, error) {
	URL := fmt.Sprintf("%s/v1/accounts:lookup?key=%s", r.identityToolkitURL, r.APIKey)

	reqBody := map[string]interface{}{
		"idToken": token,
	}

	reqBodyJSON, _ := json.Marshal(reqBody)

	resp, err := http.Post(URL, "", bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	respBuffer, _ := io.ReadAll(resp.Body)

	var userData UserDataFromToken
	err = json.Unmarshal(respBuffer, &userData)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	if userData.Kind == "" {
		var errorResponse ErrorFromFirebase

		err = json.Unmarshal(respBuffer, &errorResponse)
		if err != nil {
			return nil, errors.Wrap(ErrInternalServer, err.Error())
		}

		if strings.Contains(errorResponse.Error.Message, "INVALID_ID_TOKEN") {
			return nil, errors.Wrap(ErrInputValidation, "token invalid")
		}

		return nil, errors.Wrap(ErrInternalServer, errorResponse.Error.Message)
	}

	return &userData, nil
}

// NewRepo for initialize repo
func NewRepo(identityToolkitURL, secureTokenURL, APIKey string) Repo {
	return &repo{
		identityToolkitURL: identityToolkitURL,
		secureTokenURL:     secureTokenURL,
		APIKey:             APIKey,
	}
}
