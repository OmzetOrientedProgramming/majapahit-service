package booking

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetDetail(bookingID int) (*Detail, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(Detail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetItemWrapper(bookingID int) (*ItemsWrapper, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(ItemsWrapper)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetTicketPriceWrapper(bookingID int) (*TicketPriceWrapper, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(TicketPriceWrapper)
	return &ret, args.Error(1)
}

func TestService_GetDetailSuccess(t *testing.T) {
	bookingID := 1
	createdAtRow := time.Date(2021, time.Month(10), 26, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	bookingDetail := Detail{
		ID:             1,
		Date:           "27 Oktober 2021",
		StartTime:      "19:00",
		EndTime:        "20:00",
		Capacity:       10,
		Status:         1,
		TotalPriceItem: 100000.0,
		CreatedAt:      createdAtRow,
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 10000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "Jus Mangga Asyik",
				Image: "ini_link_gambar_1",
				Qty:   10,
				Price: 10000.0,
			},
			{
				Name:  "Pizza with Pinapple Large",
				Image: "ini_link_gambar_2",
				Qty:   2,
				Price: 150000.0,
			},
		},
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", bookingID).Return(bookingDetail, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)

	bookingDetailResult, err := mockService.GetDetail(bookingID)
	mockRepo.AssertExpectations(t)

	totalTicketPrice := ticketPriceWrapper.Price * float64(bookingDetail.Capacity)
	totalPrice := totalTicketPrice + bookingDetail.TotalPriceItem

	bookingDetail.TotalPriceTicket = totalTicketPrice
	bookingDetail.TotalPrice = totalPrice

	bookingDetail.Items = make([]ItemDetail, 2)
	bookingDetail.Items[0].Name = itemsWrapper.Items[0].Name
	bookingDetail.Items[0].Image = itemsWrapper.Items[0].Image
	bookingDetail.Items[0].Qty = itemsWrapper.Items[0].Qty
	bookingDetail.Items[0].Price = itemsWrapper.Items[0].Price

	bookingDetail.Items[1].Name = itemsWrapper.Items[1].Name
	bookingDetail.Items[1].Image = itemsWrapper.Items[1].Image
	bookingDetail.Items[1].Qty = itemsWrapper.Items[1].Qty
	bookingDetail.Items[1].Price = itemsWrapper.Items[1].Price

	assert.Equal(t, &bookingDetail, bookingDetailResult)
	assert.NotNil(t, bookingDetailResult)
	assert.NoError(t, err)
}

func TestService_GetDetailFailedCalledGetDetail(t *testing.T) {
	bookingID := 1
	var bookingDetail Detail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", bookingID).Return(bookingDetail, ErrInternalServerError)

	bookingDetailResult, err := mockService.GetDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, bookingDetailResult)
}

func TestService_GetDetailFailedCalledGetTicketPriceWrapper(t *testing.T) {
	bookingID := 1
	var bookingDetail Detail
	var ticketPriceWrapper TicketPriceWrapper

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", bookingID).Return(bookingDetail, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, ErrInternalServerError)

	bookingDetailResult, err := mockService.GetDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, bookingDetailResult)
}

func TestService_GetDetailFailedCalledGetItemWrapper(t *testing.T) {
	bookingID := 1
	var bookingDetail Detail
	var ticketPriceWrapper TicketPriceWrapper
	var itemsWrapper ItemsWrapper

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", bookingID).Return(bookingDetail, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, ErrInternalServerError)

	bookingDetailResult, err := mockService.GetDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, bookingDetailResult)
}

func TestService_GetDetailWrongInput(t *testing.T) {
	// Define input
	bookingID := -1

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	bookingDetail, err := mockService.GetDetail(bookingID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, bookingDetail)
}


func (m *MockRepository) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	args := m.Called(localID)
	ret := args.Get(0).([]Booking)
	return &ret, args.Error(1)
}

func TestService_GetMyBookingsOngoingSuccess(t *testing.T) {
	localID := "abc"
	myBookingsOngoing := []Booking{
		{
			ID:         1,
			PlaceID:    2,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       "2022-04-10",
			StartTime:  "08:00",
			EndTime:    "10:00",
			Status:     0,
			TotalPrice: 10000,
		}, 
		{
			ID:         2,
			PlaceID:    3,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       "2022-04-11",
			StartTime:  "09:00",
			EndTime:    "11:00",
			Status:     0,
			TotalPrice: 20000,
		},
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetMyBookingsOngoing", localID).Return(myBookingsOngoing, nil)

	myBookingsOngoingResult, err := mockService.GetMyBookingsOngoing(localID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &myBookingsOngoing, myBookingsOngoingResult)
	assert.NotNil(t, myBookingsOngoingResult)
	assert.NoError(t, err)
}

func TestService_GetMyBookingsOngoingWrongInput(t *testing.T) {
	// Define input
	localID := ""

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	myBookingsOngoing, err := mockService.GetMyBookingsOngoing(localID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, myBookingsOngoing)
}

func TestService_GetMyBookingsOngoingFailedCalledGetDetail(t *testing.T) {
	localID := "abc"
	var myBookingsOngoing []Booking

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetMyBookingsOngoing", localID).Return(myBookingsOngoing, ErrInternalServerError)

	myBookingsOngoingResult, err := mockService.GetMyBookingsOngoing(localID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, myBookingsOngoingResult)
}

func (m *MockRepository) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, error) {
	args := m.Called(params)
	ret := args.Get(0).(List)
	return &ret, args.Error(1)
}

func TestService_GetMyBookingsPreviousWithPaginationSuccess(t *testing.T) {
	// Define input and output
	myBookingsPrevious := List{
		Bookings: []Booking{
			{
				ID:         1,
				PlaceID:    2,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-10",
				StartTime:  "08:00",
				EndTime:    "10:00",
				Status:     0,
				TotalPrice: 10000,
			}, 
			{
				ID:         2,
				PlaceID:    3,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-11",
				StartTime:  "09:00",
				EndTime:    "11:00",
				Status:     0,
				TotalPrice: 20000,
			},
		},
		TotalCount: 2,
	}
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", params).Return(myBookingsPrevious, nil)

	// Test
	myBookingsPreviousResult, _, err := mockService.GetMyBookingsPreviousWithPagination(localID, params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &myBookingsPrevious, myBookingsPreviousResult)
	assert.NotNil(t, myBookingsPreviousResult)
	assert.NoError(t, err)
}

func TestService_GetMyBookingsPreviousWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	myBookingsPrevious := List{
		Bookings: []Booking{
			{
				ID:         1,
				PlaceID:    2,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-10",
				StartTime:  "08:00",
				EndTime:    "10:00",
				Status:     0,
				TotalPrice: 10000,
			}, 
			{
				ID:         2,
				PlaceID:    3,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-11",
				StartTime:  "09:00",
				EndTime:    "11:00",
				Status:     0,
				TotalPrice: 20000,
			},
		},
		TotalCount: 2,
	}
	localID := "abc"
	params := BookingsListRequest{
		Limit: 0,
		Page:  0,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	paramsDefault := BookingsListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", paramsDefault).Return(myBookingsPrevious, nil)

	// Test
	myBookingsPreviousResult, _, err := mockService.GetMyBookingsPreviousWithPagination(localID, params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &myBookingsPrevious, myBookingsPreviousResult)
	assert.NotNil(t, myBookingsPreviousResult)
	assert.NoError(t, err)
}

func TestService_GetMyBookingsPreviousWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	localID := "abc"
	params := BookingsListRequest{
		Limit: 101,
		Page:  0,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	myBookingsPreviousResult, _, err := mockService.GetMyBookingsPreviousWithPagination(localID, params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, myBookingsPreviousResult)
}

func TestService_GetMyBookingsPreviousWithPaginationFailedCalledGetPlacesListWithPagination(t *testing.T) {
	// Define input and output
	var myBookingsPrevious List

	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", params).Return(myBookingsPrevious, ErrInternalServerError)

	// Test
	myBookingsPreviousResult, _, err := mockService.GetMyBookingsPreviousWithPagination(localID, params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, myBookingsPreviousResult)
}
