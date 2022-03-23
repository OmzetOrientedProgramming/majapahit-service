package customerbooking

import (
	"testing"

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
