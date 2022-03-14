package businessadminauth

import (
	"strings"

	"github.com/pkg/errors"
)

// NewService is a constructor to get a Service instance
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Service is used to define the methods in it
type Service interface {
	RegisterBusinessAdmin(request RegisterBusinessAdminRequest) (*LoginCredential, error)
}

// service is a struct of service
type service struct {
	repo Repo
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
