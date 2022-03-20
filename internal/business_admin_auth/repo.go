package businessadminauth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"math/big"
	"net/mail"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo is a constructor to get Repo instance
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

// repo is a struct to define repo
type repo struct {
	db *sqlx.DB
}

// Repo is an interface to define methods in it
type Repo interface {
	GetBusinessAdminByEmail(email string) (*BusinessAdmin, error)
	CreateUser(phoneNumber, name, email, password string, status int) error                     //status = 1
	CreateBusinessAdmin(userID int, bankAccount, bankAccountName string, balance float32) error //balance = 0.0
	CreatePlace(name, address string, capacity int, description string,
		userID, interval int, openHour, closeHour, image string,
		minHourBooking, maxHourBooking, minSlotBooking, maxSlotBooking int,
		lat, long float64) error

	CheckRequiredFields(request RegisterBusinessAdminRequest, errorList []string) []string
	CheckBusinessAdminFields(request RegisterBusinessAdminRequest) error
	CheckUserFields(request RegisterBusinessAdminRequest) error
	CheckPlaceFields(request RegisterBusinessAdminRequest) error

	CheckIfPhoneNumberIsUnique(phoneNumber string) (bool, error)
	RetrieveUserID(phoneNumber string) (int, error)
	CheckIfBankAccountIsUnique(bankAccount string) (bool, error)
	CheckIfEmailIsUnique(email string) (bool, error)
	CheckIfPlaceNameIsUnique(name string) (bool, error)
	GeneratePassword() string
	VerifyHour(hour, hourName string) (bool, error)
	CompareOpenAndCloseHour(openHour, closeHour string) (bool, error)
}

// CreateUser is a method in which we insert a new row of users
func (r repo) CreateUser(phoneNumber, name, email, password string, status int) error {
	_, err := r.db.Exec("INSERT INTO users (phone_number, name, email, password, status) VALUES ($1, $2, $3, $4, $5)", phoneNumber, name, email, password, status)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

// CreateBusinessAdmin is a method in which we insert a new row of business_owners
func (r repo) CreateBusinessAdmin(userID int, bankAccount, bankAccountName string, balance float32) error {

	_, err := r.db.Exec("INSERT INTO business_owners (balance, bank_account_number, bank_account_name, user_id) VALUES ($1, $2, $3, $4)",
		balance, bankAccount, bankAccountName, userID)

	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

// CreatePlace is a method in which we insert a new row of places
func (r repo) CreatePlace(name, address string, capacity int, description string,
	userID, interval int, openHour, closeHour, image string,
	minIntervalBooking, maxIntervalBooking, minSlotBooking, maxSlotBooking int,
	lat, long float64) error {

	var sqlCommand = `INSERT INTO places (
				name, address, capacity, description, user_id, interval, open_hour, close_hour, image,
				min_interval_booking, max_interval_booking, min_slot_booking, max_slot_booking, lat, long) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	_, err := r.db.Exec(sqlCommand, name, address, capacity, description, userID, interval,
		openHour, closeHour, image, minIntervalBooking, maxIntervalBooking, minSlotBooking,
		maxSlotBooking, lat, long)

	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

// CheckIfPhoneNumberIsUnique returns true if phoneNumber is not in db
func (r repo) CheckIfPhoneNumberIsUnique(phoneNumber string) (bool, error) {
	var phoneNumberResult string
	err := r.db.Get(&phoneNumberResult, "SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1", phoneNumber)
	switch err {
	case nil:
		return false, nil
	case sql.ErrNoRows:
		return true, nil
	}
	return false, errors.Wrap(ErrInternalServerError, err.Error())
}

// RetrieveUserId is called to retrieve a new generated user_id
func (r repo) RetrieveUserID(phoneNumber string) (int, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE phone_number=$1 LIMIT 1", phoneNumber)
	if err != nil {
		return -1, errors.Wrap(ErrInternalServerError, err.Error())
	}
	userID := user.ID
	return userID, nil
}

// CheckIfBankAccountIsUnique returns true if bankAccount is not in db
func (r repo) CheckIfBankAccountIsUnique(bankAccount string) (bool, error) {
	var bankAccountResult string
	err := r.db.Get(&bankAccountResult, "SELECT bank_account_number FROM business_owners WHERE bank_account_number=$1 LIMIT 1", bankAccount)
	switch err {
	case nil:
		return false, nil
	case sql.ErrNoRows:
		return true, nil
	}
	return false, errors.Wrap(ErrInternalServerError, err.Error())
}

// CheckIfEmailIsUnique returns true if email is not in db
func (r repo) CheckIfEmailIsUnique(email string) (bool, error) {
	var emailResult string
	err := r.db.Get(&emailResult, "SELECT email FROM users WHERE email=$1 LIMIT 1", email)
	switch err {
	case nil:
		return false, nil
	case sql.ErrNoRows:
		return true, nil
	}
	return false, errors.Wrap(ErrInternalServerError, err.Error())
}

// CheckIfPlaceNameIsUnique returns true if placeName is not in db
func (r repo) CheckIfPlaceNameIsUnique(name string) (bool, error) {
	var nameResult string
	err := r.db.Get(&nameResult, "SELECT name FROM places WHERE name=$1 LIMIT 1", name)
	switch err {
	case nil:
		return false, nil
	case sql.ErrNoRows:
		return true, nil
	}
	return false, errors.Wrap(ErrInternalServerError, err.Error())
}

// CheckRequiredFields is called to ensure all required fields are fulfilled
func (r repo) CheckRequiredFields(request RegisterBusinessAdminRequest, errorList []string) []string {
	if request.AdminPhoneNumber == "" {
		errorList = append(errorList, "admin_phone_number is required")
	}

	if request.AdminEmail == "" {
		errorList = append(errorList, "admin_email is required")
	}

	if request.AdminBankAccount == "" {
		errorList = append(errorList, "admin_bank_account is required")
	}

	if request.AdminName == "" {
		errorList = append(errorList, "admin_name is required")
	}

	if request.PlaceName == "" {
		errorList = append(errorList, "place_name is required")
	}

	if request.PlaceAddress == "" {
		errorList = append(errorList, "place_address is required")
	}

	if request.PlaceCapacity == 0 {
		errorList = append(errorList, "place_capacity must be more than 0 and not empty")
	}

	if request.PlaceDescription == "" {
		errorList = append(errorList, "place_description is required")
	}

	if request.PlaceInterval == 0 {
		errorList = append(errorList, "place_interval must be more than 0 and not empty")
	}

	if request.PlaceOpenHour == "" {
		errorList = append(errorList, "place_open_hour is required")
	}

	if request.PlaceCloseHour == "" {
		errorList = append(errorList, "place_close_hour is required")
	}

	if request.PlaceImage == "" {
		errorList = append(errorList, "place_image is required")
	}

	if request.PlaceMinIntervalBooking == 0 {
		errorList = append(errorList, "place_min_interval_booking must be more than 0 and not empty")
	}

	if request.PlaceMaxIntervalBooking == 0 {
		errorList = append(errorList, "place_max_interval_booking must be more than 0 and not empty")
	}

	if request.PlaceMinSlotBooking == 0 {
		errorList = append(errorList, "place_min_slot_booking must be more than 0 and not empty")
	}

	if request.PlaceMaxSlotBooking == 0 {
		errorList = append(errorList, "place_max_slot_booking must be more than 0 and not empty")
	}

	if request.PlaceLat == 0.0 {
		errorList = append(errorList, "place_lat is required")
	}

	if request.PlaceLong == 0.0 {
		errorList = append(errorList, "place_long is required")
	}

	return errorList
}

// CheckUserFields is called to ensure all required fields regarding users are valid
func (r repo) CheckUserFields(request RegisterBusinessAdminRequest) error {
	// Check if name is valid
	if len(request.AdminName) < 3 {
		return errors.Wrap(ErrInputValidationError, "name is at least 3 characters")
	}
	if len(request.AdminName) > 50 {
		return errors.Wrap(ErrInputValidationError, "name is at most 50 characters")
	}

	// Check if phone number is valid
	for _, num := range request.AdminPhoneNumber {
		if num < '0' || num > '9' {
			return errors.Wrap(ErrInputValidationError, "phone number is invalid")
		}
	}
	if len(request.AdminPhoneNumber) > 15 {
		return errors.Wrap(ErrInputValidationError, "phone number is too long")
	}
	isPhoneNumberUnique, err := r.CheckIfPhoneNumberIsUnique(request.AdminPhoneNumber)
	if err != nil {
		return err
	}
	if !isPhoneNumberUnique {
		return errors.Wrap(ErrInputValidationError, "phone number was already taken")
	}

	// Check if user email format is valid
	_, err = mail.ParseAddress(request.AdminEmail)
	if err != nil {
		return errors.Wrap(ErrInputValidationError, "email address is invalid")
	}

	// Check if user email is unique
	isEmailUnique, err := r.CheckIfEmailIsUnique(request.AdminEmail)
	if err != nil {
		return err
	}
	if !isEmailUnique {
		return errors.Wrap(ErrInputValidationError, "email was already taken")
	}

	return nil
}

// CheckBusinessAdminFields is called to ensure all required fields regarding business_owners are valid
func (r repo) CheckBusinessAdminFields(request RegisterBusinessAdminRequest) error {
	// Check if bank_account_name is valid
	if len(request.AdminBankAccountName) < 3 {
		return errors.Wrap(ErrInputValidationError, "bank account name is at least 3 characters")
	}
	if len(request.AdminBankAccountName) > 50 {
		return errors.Wrap(ErrInputValidationError, "bank account name is at most 50 characters")
	}
	request.AdminBankAccountName = strings.ToUpper(request.AdminBankAccountName)
	for _, char := range request.AdminBankAccountName {
		if char == ' ' {
			continue
		} else if char < 'A' || char > 'Z' {
			return errors.Wrap(ErrInputValidationError, "admin bank account name is invalid")
		}
	}

	// Check if admin bank account format is valid
	if len(request.AdminBankAccount) < 10 {
		return errors.Wrap(ErrInputValidationError, "bank account is at least 10 characters")
	}
	if len(request.AdminBankAccount) > 25 {
		return errors.Wrap(ErrInputValidationError, "bank account is at most 25 characters")
	}
	for index, char := range request.AdminBankAccount {
		if index == 3 {
			if char != '-' {
				return errors.Wrap(ErrInputValidationError, "the valid bank account number format is XXX-YYY...YYY where XXX is the bank code")
			}
		} else if char < '0' || char > '9' {
			return errors.Wrap(ErrInputValidationError, "bank account number is invalid")
		}
	}

	// Check if admin bank account is unique
	isBankAccountUnique, err := r.CheckIfBankAccountIsUnique(request.AdminBankAccount)
	if err != nil {
		return err
	}
	if !isBankAccountUnique {
		return errors.Wrap(ErrInputValidationError, "a unique bank account is needed")
	}

	// Validation cleared
	return nil
}

// CheckPlaceFields is called to ensure all required fields regarding places are valid
func (r repo) CheckPlaceFields(request RegisterBusinessAdminRequest) error {
	// Place name
	if len(request.PlaceName) < 5 {
		return errors.Wrap(ErrInputValidationError, "place name is at least 5 characters")
	} else if len(request.PlaceName) > 50 {
		return errors.Wrap(ErrInputValidationError, "place name is at most 50 characters")
	}
	isPlaceNameUnique, err := r.CheckIfPlaceNameIsUnique(request.PlaceName)
	if err != nil {
		return err
	}
	if !isPlaceNameUnique {
		return errors.Wrap(ErrInputValidationError, "the place name is already taken")
	}

	// Place address
	if len(request.PlaceAddress) < 10 {
		return errors.Wrap(ErrInputValidationError, "place address is at least 10 characters")
	} else if len(request.PlaceAddress) > 100 {
		return errors.Wrap(ErrInputValidationError, "place address is at most 100 characters")
	}

	// Capacity
	if request.PlaceCapacity < 1 {
		return errors.Wrap(ErrInputValidationError, "place capacity is at least 1")
	}

	// Description
	if len(request.PlaceDescription) < 20 {
		return errors.Wrap(ErrInputValidationError, "place description is at least 20 characters")
	} else if len(request.PlaceDescription) > 2000 {
		return errors.Wrap(ErrInputValidationError, "place description is at most 2000 characters")
	}

	// Interval
	if request.PlaceInterval < 30 {
		return errors.Wrap(ErrInputValidationError, "place interval is at least 30 minutes")
	}
	if request.PlaceInterval%30 != 0 {
		return errors.Wrap(ErrInputValidationError, "place interval must be able to be divided by 30")
	}

	// OpenHour, CloseHour
	isOpenHourOkay, err := r.VerifyHour(request.PlaceOpenHour, "open hour")
	if err != nil {
		return err
	}

	isCloseHourOkay, err := r.VerifyHour(request.PlaceCloseHour, "close hour")
	if err != nil {
		return err
	}

	if isOpenHourOkay && isCloseHourOkay {
		isHourOkay, err := r.CompareOpenAndCloseHour(request.PlaceOpenHour, request.PlaceCloseHour)
		if err != nil {
			return err
		} else if !isHourOkay {
			return errors.Wrap(ErrInputValidationError, "open hour procedes close hour")
		}
	}

	// Image
	placeImageSize := len(request.PlaceImage)
	if request.PlaceImage[:32] != "https://drive.google.com/file/d/" ||
		request.PlaceImage[(placeImageSize-17):] != "/view?usp=sharing" {
		return errors.Wrap(ErrInputValidationError, "image has to be formatted as https://drive.google.com/file/d/.../view?usp=sharing")
	}

	// minIntervalBooking, maxIntervalBooking
	if request.PlaceMinIntervalBooking < 1 {
		return errors.Wrap(ErrInputValidationError, "min interval booking must be at least 1")
	} else if request.PlaceMinIntervalBooking > request.PlaceMaxIntervalBooking {
		return errors.Wrap(ErrInputValidationError, "min interval booking is more than max interval booking")
	}

	// minSlotBooking, maxSlotBooking
	if request.PlaceMinSlotBooking < 1 {
		return errors.Wrap(ErrInputValidationError, "min slot booking must be at least 1")
	} else if request.PlaceMinSlotBooking > request.PlaceMaxSlotBooking {
		return errors.Wrap(ErrInputValidationError, "min slot booking is more than max slot booking")
	}

	// lat, long
	if request.PlaceLat < 94.5 || request.PlaceLat > 141.5 {
		return errors.Wrap(ErrInputValidationError, "latitude of the place is out of reach")
	} else if request.PlaceLong < -11.5 || request.PlaceLong > 6.5 {
		return errors.Wrap(ErrInputValidationError, "longitude of the place is out of reach")
	}

	return nil
}

// VerifyHour is called to verify whether hour is in the right format
func (r repo) VerifyHour(hour, hourName string) (bool, error) {
	if len(hour) != 5 {
		return false, errors.Wrap(ErrInputValidationError, fmt.Sprintf("please use HH:mm format for %s", hourName))
	}

	for index, char := range hour {
		if index == 2 && char != ':' {
			return false, errors.Wrap(ErrInputValidationError, fmt.Sprintf("please use HH:mm format for %s", hourName))
		} else if index == 2 && char == ':' {
			continue
		} else if char < '0' || char > '9' {
			return false, errors.Wrap(ErrInputValidationError, fmt.Sprintf("please use HH:mm format for %s", hourName))
		}
	}

	hourHourtime, err := strconv.ParseInt(hour[0:1], 10, 64)
	if err != nil {
		return false, errors.Wrap(ErrInternalServerError, "error while parsing hour")
	}

	hourMinutetime, err := strconv.ParseInt(hour[3:4], 10, 64)

	if err != nil {
		return false, errors.Wrap(ErrInternalServerError, "error while parsing hour")
	}
	if hourHourtime > 23 {
		return false, errors.Wrap(ErrInputValidationError, fmt.Sprintf("hour time for %s is too large", hourName))
	}
	if hourMinutetime > 59 {
		return false, errors.Wrap(ErrInputValidationError, fmt.Sprintf("minute time for %s is too large", hourName))
	}

	return true, nil
}

// CompareOpenAndCloseHour returns true if openHour procedes closeHour
func (r repo) CompareOpenAndCloseHour(openHour, closeHour string) (bool, error) {
	openHourHourtime, err := strconv.ParseInt(openHour[0:2], 10, 64)
	if err != nil {
		return false, err
	}
	openHourMinutetime, err := strconv.ParseInt(openHour[3:5], 10, 64)
	if err != nil {
		return false, err
	}
	closeHourHourtime, err := strconv.ParseInt(closeHour[0:2], 10, 64)
	if err != nil {
		return false, err
	}
	closeHourMinutetime, err := strconv.ParseInt(closeHour[3:5], 10, 64)
	if err != nil {
		return false, err
	}

	if openHourHourtime > closeHourHourtime {
		return false, nil
	} else if openHourMinutetime > closeHourMinutetime {
		return false, nil
	} else {
		return true, nil
	}
}

// GeneratePassword returns a newly generated password
func (r repo) GeneratePassword() string {
	const characters string = "abcdefghijklmnopqrstuvwxyz0123456789"
	var password = ""

	for i := 0; i < 8; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(36))
		if err != nil {
			return ""
		}
		password = password + string(characters[idx.Int64()])
	}

	return password
}

// GetBusinessAdminByEmail returns business admin by given email
func (r repo) GetBusinessAdminByEmail(email string) (*BusinessAdmin, error) {
	var userModel auth.UserModel
	query := "SELECT * FROM users WHERE email = $1"
	if err := r.db.Get(&userModel, query, email); err != nil {
		logrus.Error("[error while accessing user model]", err.Error())
		return nil, fmt.Errorf("cannot access user table: %w", ErrInternalServerError)
	}

	var businessAdminModel BusinessAdminModel
	query = "SELECT * FROM business_owners WHERE user_id = $1"
	if err := r.db.Get(&businessAdminModel, query, userModel.ID); err != nil {
		logrus.Error("[error while accessing user model]", err.Error())
		return nil, fmt.Errorf("cannot access business owner table: %w", ErrInternalServerError)
	}

	return &BusinessAdmin{
		ID:                userModel.ID,
		Name:              userModel.Name,
		PhoneNumber:       userModel.PhoneNumber,
		Email:             userModel.Email,
		Password:          userModel.Password,
		Status:            userModel.Status,
		Balance:           businessAdminModel.Balance,
		BankAccountNumber: businessAdminModel.BankAccountNumber,
		BankAccountName:   businessAdminModel.BankAccountName,
	}, nil
}
