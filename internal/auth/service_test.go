package auth

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestService_CheckPhoneNumber(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, TwillioCredentials{
			SID:        "mockSID",
			AccountSID: "mockAccountSID",
			AuthToken:  "mockAuthToken",
		})
		mockRepo.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		actual, err := mockService.CheckPhoneNumber("087748176534")
		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.True(t, actual)
	})

	t.Run("phone number is not valid", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, TwillioCredentials{
			SID:        "mockSID",
			AccountSID: "mockAccountSID",
			AuthToken:  "mockAuthToken",
		})
		actual, err := mockService.CheckPhoneNumber("abcdefg")
		assert.Equal(t, errors.Cause(err), ErrInputValidation)
		assert.NotNil(t, actual)
		assert.False(t, actual)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, TwillioCredentials{
			SID:        "mockSID",
			AccountSID: "mockAccountSID",
			AuthToken:  "mockAuthToken",
		})
		mockRepo.On("CheckPhoneNumber", "087748176534").Return(false, ErrInternalServer)
		actual, err := mockService.CheckPhoneNumber("087748176534")
		mockRepo.AssertExpectations(t)
		assert.Equal(t, errors.Cause(err), ErrInternalServer)
		assert.NotNil(t, actual)
		assert.False(t, actual)
	})
}

func TestService_SendOTP(t *testing.T) {
	//t.Run("success", func(t *testing.T) {
	//	mockRepo := new(MockRepository)
	//	mockService := NewService(mockRepo, TwillioCredentials{
	//		SID:        "VA844307e9628fb83a3e56b18790e40c66",
	//		AccountSID: "AC5c5eabf384261619dca4c8d9eb2480dd",
	//		AuthToken:  "49a2eec79210b79e9afb82a9dc5b2f4f",
	//	})
	//	assert.NoError(t, mockService.SendOTP("085156934378"))
	//})
}

func TestService_VerifyOTP(t *testing.T) {
	//t.Run("success", func(t *testing.T) {
	//	mockRepo := new(MockRepository)
	//	mockService := NewService(mockRepo, TwillioCredentials{
	//		SID:        "mockSID",
	//		AccountSID: "mockAccountSID",
	//		AuthToken:  "mockAuthToken",
	//	})
	//	status, err := mockService.VerifyOTP("085156934378", "random")
	//	assert.NoError(t, err)
	//	assert.True(t, status)
	//})
}

func TestService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, TwillioCredentials{
			SID:        "mockSID",
			AccountSID: "mockAccountSID",
			AuthToken:  "mockAuthToken",
		})
		mockRepo.On("CreateCustomer", Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      1,
		}).Return(&Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      1,
		}, nil)

		actual, err := mockService.Register(Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      1,
		})
		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, *actual, Customer{
			PhoneNumber: "081223784562",
			Name:        "Bambang",
			Status:      1,
		})
	})
}
