package customer

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) PutEditCustomer(customer EditCustomerRequest) error {
	args := m.Called(customer)
	return args.Error(0)
}

func TestService_PutEditCustomer(t *testing.T) {
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	t.Run("Put edit customer done successfully", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		mockRepo.On("PutEditCustomer", customer).Return(nil)
		err := mockService.PutEditCustomer(customer)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Empty name input validation", func(t *testing.T) {
		userID := 1
		name := ""
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Name less than three input validation", func(t *testing.T) {
		userID := 1
		name := "yo"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Profile picture empty input validation", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := ""
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Date of birth invalid input validation", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "invalid date"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Gender invalid input validation", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "invalid date"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 3

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Repo error handling", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-09-04"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 2

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		mockRepo.On("PutEditCustomer", customer).Return(errors.Wrap(ErrInputValidation, "error repo"))
		err := mockService.PutEditCustomer(customer)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})
}
