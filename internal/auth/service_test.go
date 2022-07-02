package auth

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CheckPhoneNumber(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) CreateCustomer(customer Customer) (*Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(*Customer), args.Error(1)
}

func (m *MockRepository) GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error) {
	args := m.Called(phoneNumber)
	return args.Get(0).(*Customer), args.Error(1)
}

type FirebaseMockRepository struct {
	mock.Mock
}

func (f *FirebaseMockRepository) SendOTP(params firebaseauth.SendOTPParams) (*firebaseauth.SendOTPResult, error) {
	args := f.Called(params)
	return args.Get(0).(*firebaseauth.SendOTPResult), args.Error(1)
}

func (f *FirebaseMockRepository) VerifyOTP(params firebaseauth.VerifyOTPParams) (*firebaseauth.VerifyOTPResult, error) {
	args := f.Called(params)
	return args.Get(0).(*firebaseauth.VerifyOTPResult), args.Error(1)
}

func (f *FirebaseMockRepository) GetUserDataFromToken(token string) (*firebaseauth.UserDataFromToken, error) {
	args := f.Called(token)
	return args.Get(0).(*firebaseauth.UserDataFromToken), args.Error(1)
}

func TestService_CheckPhoneNumber(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CheckPhoneNumber", "+6287748176534").Return(true, nil)
		actual, err := mockService.CheckPhoneNumber("087748176534")
		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.True(t, actual)
	})

	t.Run("phone number is not valid", func(t *testing.T) {
		actual, err := mockService.CheckPhoneNumber("abcdefg")
		assert.Equal(t, errors.Cause(err), ErrInputValidation)
		assert.NotNil(t, actual)
		assert.False(t, actual)
	})
}

func TestService_CheckPhoneNumberErrorInternalServer(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("CheckPhoneNumber", "+6287748176534").Return(false, ErrInternalServer)
		actual, err := mockService.CheckPhoneNumber("087748176534")
		mockRepo.AssertExpectations(t)
		assert.Equal(t, errors.Cause(err), ErrInternalServer)
		assert.NotNil(t, actual)
		assert.False(t, actual)
	})
}

func TestService_SendOTPFirebaseErrorInternalServer(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("internal server error", func(t *testing.T) {
		otpResult := firebaseauth.SendOTPResult{}
		firebaseParams := firebaseauth.SendOTPParams{
			PhoneNumber:    "+628121212121",
			RecaptchaToken: "test token",
		}

		firebaseMockRepo.On("SendOTP", firebaseParams).Return(&otpResult, errors.Wrap(ErrInternalServer, "test internal server error"))
		res, err := mockService.SendOTP("08121212121", firebaseParams.RecaptchaToken)
		firebaseMockRepo.AssertExpectations(t)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Equal(t, "", res)
	})
}

func TestService_SendOTPFirebaseErrorInputValidation(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("input validation error from firebase", func(t *testing.T) {
		otpResult := firebaseauth.SendOTPResult{}
		firebaseParams := firebaseauth.SendOTPParams{
			PhoneNumber:    "+628121212121",
			RecaptchaToken: "test token",
		}

		firebaseMockRepo.On("SendOTP", firebaseParams).Return(&otpResult, errors.Wrap(firebaseauth.ErrInputValidation, "test error input validation from firebase"))
		res, err := mockService.SendOTP("08121212121", firebaseParams.RecaptchaToken)
		firebaseMockRepo.AssertExpectations(t)

		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Equal(t, "", res)
	})
}

func TestService_SendOTP(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("success", func(t *testing.T) {
		otpResult := firebaseauth.SendOTPResult{}

		firebaseParams := firebaseauth.SendOTPParams{
			PhoneNumber:    "+628121212121",
			RecaptchaToken: "test token",
		}
		firebaseMockRepo.On("SendOTP", firebaseParams).Return(&otpResult, nil)
		res, err := mockService.SendOTP("08121212121", firebaseParams.RecaptchaToken)
		firebaseMockRepo.AssertExpectations(t)

		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("input validation error", func(t *testing.T) {
		res, err := mockService.SendOTP("", "")
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Equal(t, "", res)
	})
}

func TestService_VerifyOTP(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	var verifyOTPResult firebaseauth.VerifyOTPResult

	t.Run("success", func(t *testing.T) {
		params := firebaseauth.VerifyOTPParams{
			SessionInfo: "test session info",
			Code:        "111111",
		}

		firebaseMockRepo.On("VerifyOTP", params).Return(&verifyOTPResult, nil)
		status, err := mockService.VerifyOTP(params.SessionInfo, params.Code)
		firebaseMockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, status)
	})

	t.Run("failed input validation", func(t *testing.T) {
		status, err := mockService.VerifyOTP("", "")
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, status)
	})
}

func TestService_VerifyOTPFirebaseInputValidationError(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	var verifyOTPResult firebaseauth.VerifyOTPResult

	t.Run("failed input validation error from firebase", func(t *testing.T) {
		params := firebaseauth.VerifyOTPParams{
			SessionInfo: "test session info",
			Code:        "111111",
		}

		firebaseMockRepo.On("VerifyOTP", params).Return(&verifyOTPResult, errors.Wrap(firebaseauth.ErrInputValidation, "test firebase input validation error"))
		status, err := mockService.VerifyOTP(params.SessionInfo, params.Code)
		firebaseMockRepo.AssertExpectations(t)

		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, status)
	})
}

func TestService_VerifyOTPFirebaseInternalError(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	var verifyOTPResult firebaseauth.VerifyOTPResult

	t.Run("failed internal error", func(t *testing.T) {
		params := firebaseauth.VerifyOTPParams{
			SessionInfo: "test session info",
			Code:        "111111",
		}

		firebaseMockRepo.On("VerifyOTP", params).Return(&verifyOTPResult, errors.Wrap(firebaseauth.ErrInternalServer, "test firebase internal server error"))
		status, err := mockService.VerifyOTP(params.SessionInfo, params.Code)
		firebaseMockRepo.AssertExpectations(t)

		assert.Equal(t, firebaseauth.ErrInternalServer, errors.Cause(err))
		assert.Nil(t, status)
	})
}

func TestService_Register(t *testing.T) {
	mockRepo := new(MockRepository)
	firebaseMockRepo := new(FirebaseMockRepository)
	mockService := NewService(mockRepo, firebaseMockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreateCustomer", Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      "customer",
			Gender:      "undefined",
		}).Return(&Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      "customer",
			Gender:      "undefined",
		}, nil)

		actual, err := mockService.Register(Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      "customer",
			Gender:      "undefined",
		})
		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, *actual, Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      "customer",
			Gender:      "undefined",
		})
	})

	t.Run("input validation error", func(t *testing.T) {
		actual, err := mockService.Register(Customer{
			PhoneNumber: "invalid phone number",
			Name:        "",
			Status:      "customer",
			Gender:      "undefined",
		})

		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Nil(t, actual)
	})

	t.Run("repo error", func(t *testing.T) {
		var customer Customer
		params := Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      "customer",
			Gender:      "undefined",
			LocalID:     "test local id",
		}
		mockRepo.On("CreateCustomer", params).Return(&customer, errors.Wrap(ErrInternalServer, "test error internal"))

		actual, err := mockService.Register(params)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, actual)
	})
}

func TestService_GetCustomerByPhoneNumber(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, firebaseauth.NewRepo("test", "test", "test"))

		actual, err := mockService.GetCustomerByPhoneNumber("abc")
		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	t.Run("success", func(t *testing.T) {
		customer := Customer{}
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, firebaseauth.NewRepo("test", "test", "test"))

		phoneNumber := "0812121212121"
		mockRepo.On("GetCustomerByPhoneNumber", phoneNumber).Return(&customer, nil)
		actual, err := mockService.GetCustomerByPhoneNumber(phoneNumber)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, actual)
		assert.Nil(t, err)
		assert.Equal(t, &customer, actual)
	})
}
