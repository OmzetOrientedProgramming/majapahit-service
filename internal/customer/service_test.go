package customer

import (
	// "fmt"
	"testing"
	"time"

	// "github.com/pkg/errors"

	// "testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) RetrieveCustomerProfile(userID int) (*Profile, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func TestService_RetrieveCustomerProfile(t *testing.T) {
	mockUserID := 1

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockExpectedCustomerProfile := &Profile{
			PhoneNumber: 		"08123456789",
			Name:               "test_name_profile",
			Gender: 			0,
			DateOfBirth: 		time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
			ProfilePicture:     "test_image_profile",
		}

		mockRepo.On("RetrieveCustomerProfile", mockUserID).Return(mockExpectedCustomerProfile, nil)

		actualCustomerProfile, err := mockService.RetrieveCustomerProfile(mockUserID)
		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, actualCustomerProfile)
		assert.Equal(t, mockExpectedCustomerProfile, actualCustomerProfile)
	})

	t.Run("failed to retrieve customer profile", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockExpectedCustomerProfile := &Profile{}
		mockRepo.On("RetrieveCustomerProfile", mockUserID).Return(mockExpectedCustomerProfile, ErrInternalServerError)

		actualCustomerProfile, err := mockService.RetrieveCustomerProfile(mockUserID)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.Nil(t, actualCustomerProfile)
	})
}

