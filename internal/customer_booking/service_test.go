package customerbooking

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetListCustomerBookingWithPagination(params ListRequest) (*List, error) {
	args := m.Called(params)
	ret := args.Get(0).(List)
	return &ret, args.Error(1)
}

func TestService_GetListCustomerBookingWithPaginationSuccess(t *testing.T) {
	// Define input and output
	listCustomerBookingExpected := List{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date:         "test date 1",
				StartTime:    "test start time 1",
				EndTime:      "test end time 1",
			},
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         "test date 2",
				StartTime:    "test start time 2",
				EndTime:      "test end time 2",
			},
		},
		TotalCount: 10,
	}

	params := ListRequest{
		Limit:   10,
		Page:    1,
		Path:    "api/v1/testing",
		State:   1,
		PlaceID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetListCustomerBookingWithPagination", params).Return(listCustomerBookingExpected, nil)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}

func TestService_GetListCustomerBookingWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	listCustomerBookingExpected := List{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date:         "test date 1",
				StartTime:    "test start time 1",
				EndTime:      "test end time 1",
			},
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         "test date 2",
				StartTime:    "test start time 2",
				EndTime:      "test end time 2",
			},
		},
		TotalCount: 10,
	}

	params := ListRequest{
		Limit:   0,
		Page:    0,
		Path:    "api/v1/testing",
		State:   0,
		PlaceID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	paramsDefault := ListRequest{
		Limit:   10,
		Page:    1,
		Path:    "api/v1/testing",
		State:   1,
		PlaceID: 1,
	}

	// Expectation
	mockRepo.On("GetListCustomerBookingWithPagination", paramsDefault).Return(listCustomerBookingExpected, nil)

	// Test
	listCustomerBookingReturn, _, err := mockService.GetListCustomerBookingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listCustomerBookingExpected, listCustomerBookingReturn)
	assert.NotNil(t, listCustomerBookingReturn)
	assert.NoError(t, err)
}

func TestService_GetListCustomerBookingWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	params := ListRequest{
		Limit: 101,
		Page:  0,
		Path:  "/api/testing",
		State: 0,
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
}

func TestService_GetListCustomerBookingWithPaginationError(t *testing.T) {
	listCustomerBooking := List{}

	params := ListRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/testing",
		State:	1,
		PlaceID: 1,
	}

	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetListCustomerBookingWithPagination", params).Return(listCustomerBooking, ErrInternalServerError)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
}

func TestService_GetListCustomerBookingWithPaginationFailedURLIsEmpty(t *testing.T) {
	// Define input
	params := ListRequest{
		Limit: 100,
		Page:  0,
		Path:  "",
		State: 0,
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
}