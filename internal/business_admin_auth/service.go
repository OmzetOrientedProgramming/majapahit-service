package businessadminauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/mail"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"
)

// NewService is a constructor to get a Service instance
func NewService(repo Repo, firebaseAPIKey, identityToolkitURL string) Service {
	return &service{
		repo:               repo,
		firebaseAPIKey:     firebaseAPIKey,
		identityToolkitURL: identityToolkitURL,
	}
}

// Service is used to define the methods in it
type Service interface {
	RegisterBusinessAdmin(request RegisterBusinessAdminRequest) (*LoginCredential, error)
	Login(email, password, recaptchaToken string) (string, string, error)
}

// service is a struct of service
type service struct {
	repo               Repo
	firebaseAPIKey     string
	identityToolkitURL string
}

// RegisterBusinessAdmin is called to make users, business_owners and places, with a return of LoginCredential
func (s service) RegisterBusinessAdmin(request RegisterBusinessAdminRequest) (*LoginCredential, error) {
	var (
		credentials LoginCredential
		errorList   []string
	)

	// Check required overall
	errorList = s.repo.CheckRequiredFields(request, errorList)
	if len(errorList) != 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	// Check Fields for User
	err := s.repo.CheckUserFields(request)
	if err != nil {
		return nil, err
	}

	// Check Fields for BusinessAdmin
	err = s.repo.CheckBusinessAdminFields(request)
	if err != nil {
		return nil, err
	}

	// Check Fields for Place
	err = s.repo.CheckPlaceFields(request)
	if err != nil {
		return nil, err
	}

	// Generating Password for User
	password := s.repo.GeneratePassword()
	if password == "" {
		return nil, errors.Wrap(ErrInternalServerError, "error while generating password")
	}

	// Creating User
	var status = 1 // business_owners
	err = s.repo.CreateUser(request.AdminPhoneNumber, request.AdminName, request.AdminEmail, password, status)
	if err != nil {
		return nil, err
	}

	// Retrieving User ID
	userID, err := s.repo.RetrieveUserID(request.AdminPhoneNumber)
	if err != nil {
		return nil, err
	}

	// Creating new BusinessAdmin
	var balance float32 = 0.0
	err = s.repo.CreateBusinessAdmin(userID, request.AdminBankAccount, request.AdminBankAccountName, balance)
	if err != nil {
		return nil, err
	}

	// Creating new Place
	err = s.repo.CreatePlace(request.PlaceName, request.PlaceAddress, request.PlaceCapacity,
		request.PlaceDescription, userID, request.PlaceInterval, request.PlaceOpenHour, request.PlaceCloseHour,
		request.PlaceImage, request.PlaceMinIntervalBooking, request.PlaceMaxIntervalBooking, request.PlaceMinSlotBooking,
		request.PlaceMaxSlotBooking, request.PlaceLat, request.PlaceLong)
	if err != nil {
		return nil, err
	}

	// Inserting Login Credentials
	credentials.PlaceName = request.PlaceName
	credentials.Email = request.AdminEmail
	credentials.Password = password

	return &credentials, nil
}

// Login as a business admin
func (s service) Login(email, password, recaptchaToken string) (string, string, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return "", "", errors.Wrap(ErrInputValidationError, "invalid email address")
	}

	if password == "" {
		return "", "", errors.Wrap(ErrInputValidationError, "password cannot be empty")
	}

	businessAdmin, err := s.repo.GetBusinessAdminByEmail(email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(businessAdmin.Password), []byte(password)); err != nil {
		logrus.Errorf("[error while comparing password] %v", err)
		return "", "", errors.Wrap(ErrUnauthorized, "wrong email/password")
	}

	apiURL := fmt.Sprintf("%s/v1/accounts:signInWithPassword?key=%s", s.identityToolkitURL, s.firebaseAPIKey)
	data := map[string]interface{}{
		"email":             email,
		"password":          password,
		"captchaResponse":   recaptchaToken,
		"returnSecureToken": true,
	}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(apiURL, echo.MIMEApplicationJSON, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Error("[error while creating request] ", err.Error())
		return "", "", errors.Wrap(ErrInternalServerError, err.Error())
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		logrus.Error("[non ok status response] ", resp.StatusCode)
		return "", "", errors.Wrap(ErrInternalServerError, "non-ok status code")
	}

	var jsonResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&jsonResponse)

	return jsonResponse["idToken"].(string), jsonResponse["refreshToken"].(string), nil
}
