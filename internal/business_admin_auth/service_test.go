package businessadminauth

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CheckRequiredFields(request RegisterBusinessAdminRequest, errorList []string) []string {
	args := m.Called(request, errorList)
	ret := args.Get(0).([]string)
	return ret
}

func (m *MockRepository) CheckUserFields(request RegisterBusinessAdminRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockRepository) GeneratePassword() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRepository) CreateUser(phoneNumber, name, email, password string, status int) error {
	args := m.Called(phoneNumber, name, email, password, status)
	return args.Error(0)
}

func (m *MockRepository) CheckBusinessAdminFields(request RegisterBusinessAdminRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockRepository) CreateBusinessAdmin(userID int, bankAccount, bankAccountName string, balance float32) error {
	args := m.Called(userID, bankAccount, bankAccountName, balance)
	return args.Error(0)
}

func (m *MockRepository) CheckPlaceFields(request RegisterBusinessAdminRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockRepository) CreatePlace(name, address string, capacity int, description string,
	userID, interval int, openHour, closeHour, image string,
	minHourBooking, maxHourBooking, minSlotBooking, maxSlotBooking int,
	lat, long float64) error {
	args := m.Called(name, address, capacity, description, userID, interval, openHour,
		closeHour, image, minHourBooking, maxHourBooking, minSlotBooking, maxSlotBooking,
		lat, long)
	return args.Error(0)
}

func (m *MockRepository) CheckIfPhoneNumberIsUnique(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) RetrieveUserID(phoneNumber string) (int, error) {
	args := m.Called(phoneNumber)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) CheckIfBankAccountIsUnique(bankAccount string) (bool, error) {
	args := m.Called(bankAccount)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CheckIfEmailIsUnique(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CheckIfPlaceNameIsUnique(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) VerifyHour(hour, hourName string) (bool, error) {
	args := m.Called(hour, hourName)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CompareOpenAndCloseHour(openHour, closeHour string) (bool, error) {
	args := m.Called(openHour, closeHour)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetBusinessAdminByEmail(email string) (*BusinessAdmin, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BusinessAdmin), args.Error(1)
}

func TestService_RegisterBusinessAdmin(t *testing.T) {
	mockAPIKey := "1234567890"
	mockIdentityToolkitURL := "https://test"
	request := RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, mockAPIKey, mockIdentityToolkitURL)

	var mockEmptyErrorList []string
	mockRepo.On("CheckRequiredFields", request, mockEmptyErrorList).Return(mockEmptyErrorList)
	mockRepo.On("CheckUserFields", request).Return(nil)
	mockRepo.On("GeneratePassword").Return("12345678") //SKIPPED
	mockPassword := "12345678"
	mockStatus := 1
	mockRepo.On("CreateUser", request.AdminPhoneNumber, request.AdminName, request.AdminEmail,
		mockPassword, mockStatus).Return(nil)

	mockRepo.On("RetrieveUserID", request.AdminPhoneNumber).Return(1, nil)
	mockRepo.On("CheckBusinessAdminFields", request).Return(nil)

	mockUserID := 1
	var mockBalance float32 = 0.0
	mockRepo.On("CreateBusinessAdmin", mockUserID, request.AdminBankAccount, request.AdminBankAccountName, mockBalance).Return(nil)

	mockRepo.On("CheckPlaceFields", request).Return(nil)
	mockRepo.On("CreatePlace", request.PlaceName, request.PlaceAddress, request.PlaceCapacity,
		request.PlaceDescription, mockUserID, request.PlaceInterval, request.PlaceOpenHour, request.PlaceCloseHour,
		request.PlaceImage, request.PlaceMinIntervalBooking, request.PlaceMaxIntervalBooking, request.PlaceMinSlotBooking,
		request.PlaceMaxSlotBooking, request.PlaceLat, request.PlaceLong).Return(nil)

	loginCredentialResult, err := mockService.RegisterBusinessAdmin(request)
	mockRepo.AssertExpectations(t)

	assert.NoError(t, err)
	assert.NotNil(t, loginCredentialResult)
	assert.Equal(t, request.PlaceName, loginCredentialResult.PlaceName)
	assert.Equal(t, request.AdminEmail, loginCredentialResult.Email)
}

func TestService_Login(t *testing.T) {
	mockAPIKey := "1234567890"
	mockIdentityToolkitURL := "https://stoplight.io/mocks/oop-ppl/wave-api/41710619/mock/firebase"
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, mockAPIKey, mockIdentityToolkitURL)
	mockEmail := "mock@email.com"
	mockPassword := "$2a$10$PTygXSgiOxD1WxRMtIa5leSqVj7R80KY.yx9XL046UPl35ztS/cxu" // plaintext: mockpass
	mockRecaptchaToken := "asdfasdf"

	expectedBusinessAdmin := BusinessAdmin{
		ID:                1,
		Name:              "Mock Name",
		PhoneNumber:       "081223901234",
		Email:             mockEmail,
		Password:          mockPassword,
		Status:            1,
		Balance:           1000,
		BankAccountNumber: "081241234123414534",
		BankAccountName:   "BCA",
	}

	t.Run("success", func(t *testing.T) {
		expectedToken := "asdfasf"

		mockRepo.On("GetBusinessAdminByEmail", mockEmail).Return(&expectedBusinessAdmin, nil)

		actualBusinessAdmin, actualToken, err := mockService.Login(mockEmail, "mockpass", mockRecaptchaToken)
		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.Equal(t, expectedBusinessAdmin.ID, actualBusinessAdmin.ID)
		assert.Equal(t, expectedToken, actualToken)
	})

	t.Run("invalid email address", func(t *testing.T) {
		_, _, err := mockService.Login("randomwrongemailformat", mockPassword, mockRecaptchaToken)
		assert.True(t, errors.Is(err, ErrInputValidationError))
	})

	t.Run("empty password", func(t *testing.T) {
		_, _, err := mockService.Login(mockEmail, "", mockRecaptchaToken)
		assert.True(t, errors.Is(err, ErrInputValidationError))
	})

	t.Run("failed to get business admin by email", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, mockAPIKey, mockIdentityToolkitURL)
		mockRepo.On("GetBusinessAdminByEmail", mockEmail).Return(nil, ErrInternalServerError)

		_, _, err := mockService.Login(mockEmail, mockPassword, mockRecaptchaToken)
		fmt.Println(errors.Is(err, ErrInputValidationError))
		assert.True(t, errors.Is(err, ErrInternalServerError))
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo.On("GetBusinessAdminByEmail", mockEmail).Return(&BusinessAdmin{Password: "test"}, nil)
		_, _, err := mockService.Login(mockEmail, "wrong password", mockRecaptchaToken)
		assert.True(t, errors.Is(err, ErrUnauthorized))
	})

	t.Run("non ok response code", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, mockAPIKey, "https://wrongurl")

		mockRepo.On("GetBusinessAdminByEmail", mockEmail).Return(&expectedBusinessAdmin, nil)
		_, _, err := mockService.Login(mockEmail, "mockpass", mockRecaptchaToken)
		assert.True(t, errors.Is(err, ErrInternalServerError))
	})
}
