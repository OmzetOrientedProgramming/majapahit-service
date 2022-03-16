package auth

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type service struct {
	repo               Repo
	twillioCredentials TwillioCredentials
}

// NewService for initialize service
func NewService(repo Repo, twillioCredentials TwillioCredentials) Service {
	return &service{
		repo:               repo,
		twillioCredentials: twillioCredentials,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	CheckPhoneNumber(phoneNumber string) (bool, error)

	SendOTP(phoneNumber string) error

	VerifyOTP(phoneNumber, otp string) (bool, error)

	Register(customer Customer) (*Customer, error)

	GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error)
}

func (s service) CheckPhoneNumber(phoneNumber string) (bool, error) {
	//var errorList []string
	for _, num := range phoneNumber {
		if num < '0' || num > '9' {
			return false, errors.Wrap(ErrInputValidation, "phone number is invalid")
		}
	}

	exist, err := s.repo.CheckPhoneNumber(phoneNumber)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (s service) SendOTP(phoneNumber string) error {
	data := url.Values{}
	data.Set("To", "+62"+strings.TrimLeft(phoneNumber, "0"))
	data.Set("Channel", "sms")
	req, err := http.NewRequest(http.MethodPost,
		"https://verify.twilio.com/v2/Services/"+s.twillioCredentials.SID+"/Verifications",
		strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	log.Println(req)
	req.SetBasicAuth(s.twillioCredentials.AccountSID, s.twillioCredentials.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var twillioResponse map[string]interface{}
	if err := json.Unmarshal(body, &twillioResponse); err != nil {
		return err
	}
	status, ok := twillioResponse["status"].(string)
	if !ok {
		return errors.New("failed to parse twillio error code response")
	}

	if status != "pending" {
		return errors.New("failed to send OTP")
	}

	return nil
}

func (s service) VerifyOTP(phoneNumber, otp string) (bool, error) {
	data := url.Values{}
	data.Set("To", "+62"+strings.TrimLeft(phoneNumber, "0"))
	data.Set("Code", otp)

	req, err := http.NewRequest(http.MethodPost,
		"https://verify.twilio.com/v2/Services/"+s.twillioCredentials.SID+"/VerificationCheck",
		strings.NewReader(data.Encode()))
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(s.twillioCredentials.AccountSID, s.twillioCredentials.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var twillioResponse map[string]interface{}
	if err := json.Unmarshal(body, &twillioResponse); err != nil {
		return false, err
	}
	status, ok := twillioResponse["status"].(string)
	if !ok {
		return false, errors.New("failed to parse twillio error code response")
	}

	if status != "approved" {
		return false, nil
	}

	return true, nil
}

func (s service) Register(customer Customer) (*Customer, error) {
	for _, num := range customer.PhoneNumber {
		if num < '0' || num > '9' {
			return nil, errors.Wrap(ErrInputValidation, "phone number is invalid")
		}
	}

	createdCustomer, err := s.repo.CreateCustomer(Customer{
		DateOfBirth: time.Time{},
		Gender:      false,
		PhoneNumber: customer.PhoneNumber,
		Name:        customer.Name,
		Status:      1,
	})
	if err != nil {
		return nil, nil
	}

	return createdCustomer, nil
}

func (s service) GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error) {
	for _, num := range phoneNumber {
		if num < '0' || num > '9' {
			return nil, errors.Wrap(ErrInputValidation, "phone number is invalid")
		}
	}

	return s.repo.GetCustomerByPhoneNumber(phoneNumber)
}
