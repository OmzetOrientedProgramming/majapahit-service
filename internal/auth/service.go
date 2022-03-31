package auth

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"strings"
	"time"
)

type service struct {
	repo             Repo
	firebaseAuthRepo firebaseauth.Repo
}

// NewService for initialize service
func NewService(repo Repo, firebaseAuthRepo firebaseauth.Repo) Service {
	return &service{
		repo:             repo,
		firebaseAuthRepo: firebaseAuthRepo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	CheckPhoneNumber(phoneNumber string) (bool, error)
	SendOTP(phoneNumber, recaptchaToken string) (string, error)
	VerifyOTP(sessionInfo, otp string) (*VerifyOTPResult, error)
	Register(customer Customer) (*Customer, error)
	GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error)
}

func (s service) CheckPhoneNumber(phoneNumber string) (bool, error) {
	var errorList []string
	for _, num := range phoneNumber {
		if num < '0' || num > '9' {
			errorList = append(errorList, "phone number is invalid")
		}
	}

	if len(errorList) > 0 {
		return false, errors.Wrap(ErrInputValidation, strings.Join(errorList, ","))
	}

	exist, err := s.repo.CheckPhoneNumber(fmt.Sprintf("+62%s", strings.TrimLeft(phoneNumber, "0")))
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (s service) SendOTP(phoneNumber, recaptchaToken string) (string, error) {
	var errorList []string

	if phoneNumber == "" {
		errorList = append(errorList, "phone number is required")
	}

	if recaptchaToken == "" {
		errorList = append(errorList, "recaptcha token is required")
	}

	if len(errorList) > 0 {
		return "", errors.Wrap(ErrInputValidation, strings.Join(errorList, ","))
	}

	firebaseParams := firebaseauth.SendOTPParams{
		PhoneNumber:    fmt.Sprintf("+62%s", strings.TrimLeft(phoneNumber, "0")),
		RecaptchaToken: recaptchaToken,
	}

	resp, err := s.firebaseAuthRepo.SendOTP(firebaseParams)
	if err != nil {
		if errors.Cause(err) == firebaseauth.ErrInputValidation {
			return "", errors.Wrap(ErrInputValidation, err.Error())
		}
		return "", err
	}

	return resp.SessionInfo, nil
}

func (s service) VerifyOTP(sessionInfo, otp string) (*VerifyOTPResult, error) {
	var errorList []string

	if sessionInfo == "" {
		errorList = append(errorList, "session info is required")
	}

	if otp == "" {
		errorList = append(errorList, "otp code is required")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidation, strings.Join(errorList, ","))
	}

	firebaseParams := firebaseauth.VerifyOTPParams{
		SessionInfo: sessionInfo,
		Code:        otp,
	}

	resp, err := s.firebaseAuthRepo.VerifyOTP(firebaseParams)
	if err != nil {
		if errors.Cause(err) == firebaseauth.ErrInputValidation {
			return nil, errors.Wrap(ErrInputValidation, err.Error())
		}
		return nil, err
	}

	result := VerifyOTPResult{
		AccessToken:  resp.IDToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		LocalID:      resp.LocalID,
		IsNewUser:    resp.IsNewUser,
		PhoneNumber:  resp.PhoneNumber,
	}

	return &result, nil
}

func (s service) Register(customer Customer) (*Customer, error) {
	var errorList []string

	if customer.Name == "" {
		errorList = append(errorList, "name is required")
	}

	if !util.PhoneNumberRegex.MatchString(customer.PhoneNumber) {
		errorList = append(errorList, "phone number is invalid")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidation, strings.Join(errorList, ","))
	}

	createdCustomer, err := s.repo.CreateCustomer(Customer{
		DateOfBirth: time.Time{},
		Gender:      "undefined",
		PhoneNumber: customer.PhoneNumber,
		Name:        customer.Name,
		Status:      "customer",
		LocalID:     customer.LocalID,
	})
	if err != nil {
		return nil, err
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
