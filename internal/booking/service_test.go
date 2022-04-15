package booking

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	xendit2 "github.com/xendit/xendit-go"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, error) {
	args := m.Called(params)
	ret := args.Get(0).(ListBooking)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetDetail(bookingID int) (*Detail, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(Detail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	args := m.Called(localID)
	ret := args.Get(0).([]Booking)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, error) {
	args := m.Called(localID, params)
	ret := args.Get(0).(List)
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
func (m *MockRepository) GetBookingData(params GetBookingDataParams) (*[]DataForCheckAvailableSchedule, error) {
	args := m.Called(params)
	return args.Get(0).(*[]DataForCheckAvailableSchedule), args.Error(1)
}

func (m *MockRepository) GetTimeSlotsData(placeID int, selectedDate ...time.Time) (*[]TimeSlot, error) {
	args := m.Called(placeID, selectedDate)
	return args.Get(0).(*[]TimeSlot), args.Error(1)
}

func (m *MockRepository) GetPlaceCapacity(placeID int) (*PlaceOpenHourAndCapacity, error) {
	args := m.Called(placeID)
	return args.Get(0).(*PlaceOpenHourAndCapacity), args.Error(1)
}

func (m *MockRepository) CheckedItem(ids []CheckedItemParams) (*[]CheckedItemParams, bool, error) {
	args := m.Called(ids)
	return args.Get(0).(*[]CheckedItemParams), args.Bool(1), args.Error(2)
}

func (m *MockRepository) CreateBookingItems(items []CreateBookingItemsParams) (*CreateBookingItemsResponse, error) {
	args := m.Called(items)
	return args.Get(0).(*CreateBookingItemsResponse), args.Error(1)
}

func (m *MockRepository) CreateBooking(booking CreateBookingParams) (*CreateBookingResponse, error) {
	args := m.Called(booking)
	return args.Get(0).(*CreateBookingResponse), args.Error(1)
}

func (m *MockRepository) UpdateTotalPrice(params UpdateTotalPriceParams) (bool, error) {
	args := m.Called(params)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) UpdateBookingStatusByXenditID(xenditID string, status int) error {
	args := m.Called(xenditID, status)
	return args.Error(0)
}

func (m *MockRepository) GetInvoicesFromBooking(ID int) (bool, error) {
	args := m.Called(ID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) AddExpiredPayment(ID int, expiredAt time.Time) error {
	args := m.Called(ID, expiredAt)
	return args.Error(0)
}

type MockXenditService struct {
	mock.Mock
}

func (x *MockXenditService) CreateInvoice(params xendit.CreateInvoiceParams) (*xendit2.Invoice, error) {
	args := x.Called(params)
	return args.Get(0).(*xendit2.Invoice), args.Error(1)
}

func (x *MockXenditService) CreateDisbursement(params xendit.CreateDisbursementParams) (*xendit2.Disbursement, error) {
	args := x.Called(params)
	return args.Get(0).(*xendit2.Disbursement), args.Error(1)
}

func (x *MockXenditService) GetInvoice(ID string) (*xendit2.Invoice, error) {
	args := x.Called(ID)
	return args.Get(0).(*xendit2.Invoice), args.Error(1)
}

func (x *MockXenditService) GetDisbursement(ID string) (*xendit2.Disbursement, error) {
	args := x.Called(ID)
	return args.Get(0).(*xendit2.Disbursement), args.Error(1)
}

func (m *MockRepository) UpdateBookingStatus(bookingID int, newStatus int) error {
	args := m.Called(bookingID, newStatus)
	return args.Error(0)
}

func (m *MockRepository) InsertXenditInformation(params XenditInformation) (bool, error) {
	args := m.Called(params)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetPlaceBookingPrice(placeID int) (float64, error) {
	args := m.Called(placeID)
	return args.Get(0).(float64), args.Error(1)
}

func TestService_GetListCustomerBookingWithPaginationSuccess(t *testing.T) {
	// Define input and output
	date := time.Now()
	startTime := time.Now().Add(time.Duration(2 * time.Hour))
	endTime := time.Now().Add(time.Duration(3 * time.Hour))
	listCustomerBookingOutput := ListBooking{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date:         date.Add(time.Duration(-5 * time.Hour)),
				StartTime:    startTime.Add(time.Duration(-5 * time.Hour)),
				EndTime:      endTime.Add(time.Duration(-4 * time.Hour)),
			},
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         date.Add(time.Duration(2 * time.Hour)),
				StartTime:    startTime,
				EndTime:      endTime,
			},
		},
		TotalCount: 2,
	}

	listCustomerBookingExpected := ListBooking{
		CustomerBookings: []CustomerBooking{
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         date.Add(time.Duration(2 * time.Hour)),
				StartTime:    startTime,
				EndTime:      endTime,
			},
		},
		TotalCount: 1,
	}

	params := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "api/v1/testing",
		State:  1,
		UserID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Expectation
	mockRepo.On("GetListCustomerBookingWithPagination", params).Return(listCustomerBookingOutput, nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingGagal).Return(nil)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}

func TestService_GetListCustomerBookingWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	date := time.Now().Add(time.Duration(2 * time.Hour))
	startTime := time.Now().Add(time.Duration(2 * time.Hour))
	endTime := time.Now().Add(time.Duration(3 * time.Hour))
	listCustomerBookingReturned := ListBooking{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date:         time.Now(),
				StartTime:    time.Now(),
				EndTime:      time.Now(),
			},
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         date,
				StartTime:    startTime,
				EndTime:      endTime,
			},
		},
		TotalCount: 2,
	}

	listCustomerBookingExpected := ListBooking{
		CustomerBookings: []CustomerBooking{
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         date,
				StartTime:    startTime,
				EndTime:      endTime,
			},
		},
		TotalCount: 1,
	}

	params := ListRequest{
		Limit:  0,
		Page:   0,
		Path:   "api/v1/testing",
		State:  -1,
		UserID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	paramsDefault := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "api/v1/testing",
		State:  0,
		UserID: 1,
	}

	// Expectation
	mockRepo.On("GetListCustomerBookingWithPagination", paramsDefault).Return(listCustomerBookingReturned, nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingGagal).Return(nil)

	// Test
	listCustomerBookingReturn, _, err := mockService.GetListCustomerBookingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listCustomerBookingExpected, listCustomerBookingReturn)
	assert.NotNil(t, listCustomerBookingReturn)
	assert.NoError(t, err)
}

func TestService_GetListCustomerBookingWithPaginationFailedUpdateStatus(t *testing.T) {
	// Define input and output
	date := time.Now()
	startTime := time.Now().Add(time.Duration(2 * time.Hour))
	endTime := time.Now().Add(time.Duration(3 * time.Hour))
	listCustomerBookingReturned := ListBooking{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date:         time.Now(),
				StartTime:    time.Now(),
				EndTime:      time.Now(),
			},
			{
				ID:           2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date:         date,
				StartTime:    startTime,
				EndTime:      endTime,
			},
		},
		TotalCount: 2,
	}

	params := ListRequest{
		Limit:  0,
		Page:   0,
		Path:   "api/v1/testing",
		State:  -1,
		UserID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	paramsDefault := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "api/v1/testing",
		State:  0,
		UserID: 1,
	}

	// Expectation
	mockRepo.On("GetListCustomerBookingWithPagination", paramsDefault).Return(listCustomerBookingReturned, nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingGagal).Return(errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
}

func TestService_GetListCustomerBookingWithPaginationError(t *testing.T) {
	listCustomerBooking := ListBooking{}

	params := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/testing",
		State:  1,
		UserID: 1,
	}

	// Mock DB
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	listCustomerBookingResult, _, err := mockService.GetListCustomerBookingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listCustomerBookingResult)
}

func TestService_GetDetailSuccess(t *testing.T) {
	bookingID := 1
	createdAtRow := time.Date(2021, time.Month(10), 26, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	bookingDetail := Detail{
		ID:             1,
		Date:           time.Now(),
		StartTime:      time.Now(),
		EndTime:        time.Now(),
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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	mockRepo.On("GetDetail", bookingID).Return(bookingDetail, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)

	bookingDetailResult, err := mockService.GetDetail(bookingID)
	mockRepo.AssertExpectations(t)

	totalTicketPrice := ticketPriceWrapper.Price
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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	bookingDetail, err := mockService.GetDetail(bookingID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, bookingDetail)
}

func TestService_UpdateBookingStatusSuccess(t *testing.T) {
	bookingID := 1
	newStatus := 2

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)
	mockRepo.On("UpdateBookingStatus", bookingID, newStatus).Return(nil)

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Nil(t, err)
}

func TestService_UpdateBookingStatusWithBookingIDBelowOne(t *testing.T) {
	bookingID := 0
	newStatus := 1

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
}

func TestService_UpdateBookingStatusWithNewStatusBelowZero(t *testing.T) {
	bookingID := 1
	newStatus := -1

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
}

func TestService_UpdateBookingStatusFailedCalledUpdateBookingStatus(t *testing.T) {
	bookingID := 1
	newStatus := 2

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	mockRepo.On("UpdateBookingStatus", bookingID, newStatus).Return(ErrInternalServerError)

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestService_ChangeStatusToBookingBelumMembayarFailedGetInvoices(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 10000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)
	mockRepo.On("GetInvoicesFromBooking", 1).Return(false, errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Error(t, err, "test error")
}

func TestService_ChangeStatusToBookingBelumMembayarSuccess(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	xenditItems := []xendit.Item{
		{
			Name:  "test",
			Price: 100000,
			Qty:   1,
		},
	}

	invoiceParams := xendit.CreateInvoiceParams{
		PlaceID:             1,
		Items:               xenditItems,
		Description:         fmt.Sprint("Order from test"),
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		BookingFee:          20000,
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	xenditInvoiceReturned := xendit2.Invoice{
		ID:         "test id",
		InvoiceURL: "test url",
		ExpiryDate: &now,
	}

	xenditInformationParams := XenditInformation{
		XenditID:    xenditInvoiceReturned.ID,
		InvoicesURL: xenditInvoiceReturned.InvoiceURL,
		BookingID:   bookingID,
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 20000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, nil)
	mockRepo.On("GetInvoicesFromBooking", 1).Return(false, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)
	xenditService.On("CreateInvoice", invoiceParams).Return(&xenditInvoiceReturned, nil)
	mockRepo.On("InsertXenditInformation", xenditInformationParams).Return(true, nil)
	mockRepo.On("AddExpiredPayment", 1, now).Return(nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingBelumMembayar).Return(nil)

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Nil(t, err)
}

func TestService_ChangeStatusToBookingBelumMembayarFailedAddExpiredPayment(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	xenditItems := []xendit.Item{
		{
			Name:  "test",
			Price: 100000,
			Qty:   1,
		},
	}

	invoiceParams := xendit.CreateInvoiceParams{
		PlaceID:             1,
		Items:               xenditItems,
		Description:         fmt.Sprint("Order from test"),
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		BookingFee:          20000,
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	xenditInvoiceReturned := xendit2.Invoice{
		ID:         "test id",
		InvoiceURL: "test url",
		ExpiryDate: &now,
	}

	xenditInformationParams := XenditInformation{
		XenditID:    xenditInvoiceReturned.ID,
		InvoicesURL: xenditInvoiceReturned.InvoiceURL,
		BookingID:   bookingID,
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 20000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, nil)
	mockRepo.On("GetInvoicesFromBooking", 1).Return(false, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)
	xenditService.On("CreateInvoice", invoiceParams).Return(&xenditInvoiceReturned, nil)
	mockRepo.On("InsertXenditInformation", xenditInformationParams).Return(true, nil)
	mockRepo.On("AddExpiredPayment", 1, now).Return(errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Error(t, err, "test error")
}

func TestService_ChangeStatusToBookingBelumMembayarFailedInsertXenditInfo(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	xenditItems := []xendit.Item{
		{
			Name:  "test",
			Price: 100000,
			Qty:   1,
		},
	}

	invoiceParams := xendit.CreateInvoiceParams{
		PlaceID:             1,
		Items:               xenditItems,
		Description:         fmt.Sprint("Order from test"),
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		BookingFee:          20000,
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	xenditInvoiceReturned := xendit2.Invoice{
		ID:         "test id",
		InvoiceURL: "test url",
		ExpiryDate: &now,
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 20000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	xenditInformationParams := XenditInformation{
		XenditID:    xenditInvoiceReturned.ID,
		InvoicesURL: xenditInvoiceReturned.InvoiceURL,
		BookingID:   bookingID,
	}

	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, nil)
	mockRepo.On("GetInvoicesFromBooking", 1).Return(false, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)
	xenditService.On("CreateInvoice", invoiceParams).Return(&xenditInvoiceReturned, nil)
	mockRepo.On("InsertXenditInformation", xenditInformationParams).Return(false, errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Error(t, err, "test error")
}

func TestService_ChangeStatusToBookingBelumMembayarCreateInvoices(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	xenditItems := []xendit.Item{
		{
			Name:  "test",
			Price: 100000,
			Qty:   1,
		},
	}

	invoiceParams := xendit.CreateInvoiceParams{
		PlaceID:             1,
		Items:               xenditItems,
		Description:         fmt.Sprint("Order from test"),
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		BookingFee:          20000,
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	xenditInvoiceReturned := xendit2.Invoice{
		ID:         "test id",
		InvoiceURL: "test url",
		ExpiryDate: &now,
	}

	ticketPriceWrapper := TicketPriceWrapper{
		Price: 20000,
	}

	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}

	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, nil)
	mockRepo.On("GetInvoicesFromBooking", 1).Return(false, nil)
	mockRepo.On("GetTicketPriceWrapper", bookingID).Return(ticketPriceWrapper, nil)
	mockRepo.On("GetItemWrapper", bookingID).Return(itemsWrapper, nil)
	xenditService.On("CreateInvoice", invoiceParams).Return(&xenditInvoiceReturned, errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Error(t, err, "test error")
}

func TestService_ChangeStatusToBookingBelumMembayarFailedGetDetail(t *testing.T) {
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	bookingID := 1
	newStatus := util.BookingBelumMembayar

	bookingDetailOutput := Detail{
		ID:                  1,
		CustomerName:        "test",
		CustomerPhoneNumber: "test",
		PlaceID:             1,
		Date:                time.Time{},
		StartTime:           time.Time{},
		EndTime:             time.Time{},
		Capacity:            0,
		Status:              0,
		CreatedAt:           "",
		TotalPrice:          0,
		TotalPriceTicket:    0,
		TotalPriceItem:      0,
		Items: []ItemDetail{
			{
				Name:  "test",
				Image: "",
				Qty:   1,
				Price: 100000,
			},
		},
	}
	mockRepo.On("GetDetail", bookingID).Return(bookingDetailOutput, errors.Wrap(ErrInternalServerError, "test error"))

	// Test
	err := mockService.UpdateBookingStatus(bookingID, newStatus)

	assert.Error(t, err, "test error")
}

func TestService_GetMyBookingsOngoingSuccess(t *testing.T) {
	localID := "abc"
	date := time.Now().Add(time.Duration(2 * time.Hour))
	StartTime := time.Now().Add(time.Duration(2 * time.Hour))
	EndTime := time.Now().Add(time.Duration(3 * time.Hour))
	myBookingsOngoing := []Booking{
		{
			ID:         1,
			PlaceID:    2,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       time.Now(),
			StartTime:  time.Now(),
			EndTime:    time.Now(),
			Status:     0,
			TotalPrice: 10000,
		},
		{
			ID:         2,
			PlaceID:    3,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       date,
			StartTime:  StartTime,
			EndTime:    EndTime,
			Status:     0,
			TotalPrice: 20000,
		},
	}

	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	mockRepo.On("GetMyBookingsOngoing", localID).Return(myBookingsOngoing, nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingGagal).Return(nil)

	myBookingsOngoingResult, err := mockService.GetMyBookingsOngoing(localID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &myBookingsOngoing, myBookingsOngoingResult)
	assert.NotNil(t, myBookingsOngoingResult)
	assert.NoError(t, err)
}

func TestService_GetMyBookingsOngoingFailedUpdateStatus(t *testing.T) {
	localID := "abc"
	date := time.Now().Add(time.Duration(2 * time.Hour))
	StartTime := time.Now().Add(time.Duration(2 * time.Hour))
	EndTime := time.Now().Add(time.Duration(3 * time.Hour))
	myBookingsOngoing := []Booking{
		{
			ID:         1,
			PlaceID:    2,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       time.Now(),
			StartTime:  time.Now(),
			EndTime:    time.Now(),
			Status:     0,
			TotalPrice: 10000,
		},
		{
			ID:         2,
			PlaceID:    3,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       date,
			StartTime:  StartTime,
			EndTime:    EndTime,
			Status:     0,
			TotalPrice: 20000,
		},
	}

	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	mockRepo.On("GetMyBookingsOngoing", localID).Return(myBookingsOngoing, nil)
	mockRepo.On("UpdateBookingStatus", 1, util.BookingGagal).Return(errors.Wrap(ErrInternalServerError, "testerror"))

	myBookingsOngoingResult, err := mockService.GetMyBookingsOngoing(localID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, myBookingsOngoingResult)
}

func TestService_GetMyBookingsOngoingWrongInput(t *testing.T) {
	// Define input
	localID := ""

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Test
	myBookingsOngoing, err := mockService.GetMyBookingsOngoing(localID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, myBookingsOngoing)
}

func TestService_GetMyBookingsOngoingFailedCalledGetDetail(t *testing.T) {
	localID := "abc"
	var myBookingsOngoing []Booking

	mockRepo := new(MockRepository)
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	mockRepo.On("GetMyBookingsOngoing", localID).Return(myBookingsOngoing, ErrInternalServerError)

	myBookingsOngoingResult, err := mockService.GetMyBookingsOngoing(localID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, myBookingsOngoingResult)
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
				Date:       time.Now(),
				StartTime:  time.Now(),
				EndTime:    time.Now(),
				Status:     0,
				TotalPrice: 10000,
			},
			{
				ID:         2,
				PlaceID:    3,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       time.Now(),
				StartTime:  time.Now(),
				EndTime:    time.Now(),
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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", localID, params).Return(myBookingsPrevious, nil)

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
				Date:       time.Now(),
				StartTime:  time.Now(),
				EndTime:    time.Now(),
				Status:     0,
				TotalPrice: 10000,
			},
			{
				ID:         2,
				PlaceID:    3,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       time.Now(),
				StartTime:  time.Now(),
				EndTime:    time.Now(),
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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	paramsDefault := BookingsListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", localID, paramsDefault).Return(myBookingsPrevious, nil)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

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
	xenditService := new(MockXenditService)
	mockService := NewService(mockRepo, xenditService)

	// Expectation
	mockRepo.On("GetMyBookingsPreviousWithPagination", localID, params).Return(myBookingsPrevious, ErrInternalServerError)

	// Test
	myBookingsPreviousResult, _, err := mockService.GetMyBookingsPreviousWithPagination(localID, params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, myBookingsPreviousResult)
}

func TestService_makeTimeSlotsAsMap(t *testing.T) {
	concreteService := new(service)

	t.Run("success", func(t *testing.T) {
		timeSlot := []TimeSlot{
			{
				ID:        1,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Day:       1,
			},
			{
				ID:        2,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Day:       2,
			},
			{
				ID:        3,
				StartTime: time.Now(),
				EndTime:   time.Now(),
				Day:       3,
			},
		}

		expectedOutput := map[int]map[time.Time]time.Time{
			0: {},
			1: {
				timeSlot[0].StartTime: timeSlot[0].EndTime,
			},
			2: {
				timeSlot[1].StartTime: timeSlot[1].EndTime,
			},
			3: {
				timeSlot[2].StartTime: timeSlot[2].EndTime,
			},
			4: {},
			5: {},
			6: {},
		}

		timeSlotAsMap := concreteService.makeTimeSlotsAsMap(timeSlot)
		assert.Equal(t, expectedOutput, timeSlotAsMap)
	})
}

func TestService_divideBookings(t *testing.T) {
	concreteService := new(service)

	t.Run("success", func(t *testing.T) {
		fromDate, _ := time.Parse(util.DateLayout, "2022-03-29")
		checkedInterval := 2

		startDateBookingOne, _ := time.Parse(util.TimeLayout, "09:00:00")
		endDataBookingOne, _ := time.Parse(util.TimeLayout, "11:00:00")

		startDateBookingTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		endDataBookingTwo, _ := time.Parse(util.TimeLayout, "11:00:00")

		timeSlotEight, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotNine, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotTen, _ := time.Parse(util.TimeLayout, "10:00:00")
		timeSlotEleven, _ := time.Parse(util.TimeLayout, "11:00:00")

		date, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      date,
				StartTime: startDateBookingOne,
				EndTime:   endDataBookingOne,
				Capacity:  10,
			},
			{
				ID:        2,
				Date:      date,
				StartTime: startDateBookingTwo,
				EndTime:   endDataBookingTwo,
				Capacity:  10,
			},
		}

		timeSlotData := map[int]map[time.Time]time.Time{
			0: {},
			1: {},
			2: {
				timeSlotEight: timeSlotNine,
				timeSlotNine:  timeSlotTen,
				timeSlotTen:   timeSlotEleven,
			},
			3: {},
			4: {},
			5: {},
			6: {},
		}

		expectedOutput := map[string]map[string]int{
			"2022-03-29": {
				"09:00:00 - 10:00:00": 10,
				"10:00:00 - 11:00:00": 20,
			},
			"2022-03-30": {},
		}

		dividedBooking := concreteService.divideBookings(bookingData, timeSlotData, fromDate, checkedInterval)
		assert.Equal(t, expectedOutput, dividedBooking)
	})
}

func TestService_checkAvailable(t *testing.T) {
	concreteService := new(service)

	t.Run("success", func(t *testing.T) {
		placeCapacity := 20
		bookingSlot := 10

		timeSlotEight, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotNine, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotTen, _ := time.Parse(util.TimeLayout, "10:00:00")
		timeSlotEleven, _ := time.Parse(util.TimeLayout, "11:00:00")

		timeSlot := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotEight,
				EndTime:   timeSlotNine,
				Day:       2,
			},
			{
				ID:        2,
				StartTime: timeSlotNine,
				EndTime:   timeSlotTen,
				Day:       2,
			},
			{
				ID:        3,
				StartTime: timeSlotTen,
				EndTime:   timeSlotEleven,
				Day:       3,
			},
		}

		bookingData := map[string]map[string]int{
			"2022-03-29": {
				"09:00:00 - 10:00:00": 10,
				"10:00:00 - 11:00:00": 20,
			},
			"2022-03-30": {},
		}

		expectedOutput := map[string]map[string]int{
			"2022-03-29": {
				"09:00:00": 0,
				"10:00:00": 10,
			},
			"2022-03-30": {
				"11:00:00": 0,
			},
		}

		dividedBooking := concreteService.checkAvailableSchedule(bookingData, timeSlotEight, placeCapacity, bookingSlot, timeSlot, false)
		assert.Equal(t, expectedOutput, dividedBooking)
	})
}

func TestService_formatAvailableTime(t *testing.T) {
	concreteService := new(service)

	t.Run("success", func(t *testing.T) {
		date, _ := time.Parse(util.DateLayout, "2022-03-29")

		data := map[string]map[string]int{
			"2022-03-29": {
				"08:00:00 - 09:00:00": 0,
				"09:00:00 - 10:00:00": 10,
			},
			"2022-03-30": {
				"08:00:00 - 09:00:00": 0,
				"09:00:00 - 10:00:00": 0,
				"10:00:00 - 11:00:00": 0,
			},
		}

		expectedOutput := []AvailableTimeResponse{
			{
				Time:  "08:00:00 - 09:00:00",
				Total: 0,
			},
			{
				Time:  "09:00:00 - 10:00:00",
				Total: 10,
			},
		}

		dividedBooking := concreteService.formatAvailableTimeData(data, date)
		assert.Equal(t, len(expectedOutput), len(dividedBooking))
	})
}

func TestService_GetAvailableTime(t *testing.T) {
	mockRepo := new(MockRepository)
	mockXenditService := new(MockXenditService)
	mockService := NewService(mockRepo, mockXenditService)

	t.Run("success", func(t *testing.T) {
		selectedDate, _ := time.Parse(util.DateLayout, "2022-03-29")
		selectedDateSlice := []time.Time{
			selectedDate,
		}
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityOne := 10

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityTwo := 10

		params := GetAvailableTimeParams{
			PlaceID:      1,
			SelectedDate: selectedDate,
			StartTime:    startTime,
			BookedSlot:   10,
		}

		midnight := time.Date(params.SelectedDate.Year(), params.SelectedDate.Month(), params.SelectedDate.Day(), 0, 0, 0, 0, params.SelectedDate.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.SelectedDate,
			EndDate:   midnight,
			StartTime: params.StartTime,
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      selectedDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      selectedDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		timeSlotOneStartTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotTwoStartTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotThreeStartTime, _ := time.Parse(util.TimeLayout, "10:00:00")
		timeSlotFourStartTime, _ := time.Parse(util.TimeLayout, "11:00:00")

		getTimeSlotReturn := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotTwoStartTime,
				EndTime:   timeSlotThreeStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotThreeStartTime,
				EndTime:   timeSlotFourStartTime,
				Day:       2,
			},
		}

		placeCapacity := 20

		expectedOutput := []AvailableTimeResponse{
			{
				Time:  "08:00:00 - 09:00:00",
				Total: 0,
			},
			{
				Time:  "09:00:00 - 10:00:00",
				Total: 10,
			},
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", params.PlaceID, selectedDateSlice).Return(&getTimeSlotReturn, nil)
		mockRepo.On("GetPlaceCapacity", params.PlaceID).Return(&PlaceOpenHourAndCapacity{
			OpenHour: startTime,
			Capacity: placeCapacity,
		}, nil)

		output, err := mockService.GetAvailableTime(params)
		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, len(expectedOutput), len(*output))
	})

	t.Run("input validation error", func(t *testing.T) {
		mockXenditService := new(MockXenditService)
		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		params := GetAvailableTimeParams{
			PlaceID:      0,
			SelectedDate: startDate,
			StartTime:    startDate,
			BookedSlot:   0,
		}

		output, err := mockService.GetAvailableTime(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestService_GetAvailableTimeGetBookingDataFailedInternalServerError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockXenditService := new(MockXenditService)
	mockService := NewService(mockRepo, mockXenditService)

	selectedDate, _ := time.Parse(util.DateLayout, "2022-03-29")
	startTime, _ := time.Parse(util.TimeLayout, "08:00:00")

	params := GetAvailableTimeParams{
		PlaceID:      1,
		SelectedDate: selectedDate,
		StartTime:    startTime,
		BookedSlot:   10,
	}

	midnight := time.Date(params.SelectedDate.Year(), params.SelectedDate.Month(), params.SelectedDate.Day(), 0, 0, 0, 0, params.SelectedDate.Location())
	midnight = midnight.Add(time.Duration(1*24) * time.Hour)

	repoParams := GetBookingDataParams{
		PlaceID:   params.PlaceID,
		StartDate: params.SelectedDate,
		EndDate:   midnight,
		StartTime: params.StartTime,
	}

	var getBookingDataReturn []DataForCheckAvailableSchedule

	t.Run("get booking data failed", func(t *testing.T) {

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, errors.Wrap(ErrInternalServerError, "test error"))

		output, err := mockService.GetAvailableTime(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestService_GetAvailableTimeGetTimeSlotFailedInternalServerError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockXenditService := new(MockXenditService)
	mockService := NewService(mockRepo, mockXenditService)

	selectedDate, _ := time.Parse(util.DateLayout, "2022-03-29")
	startTime, _ := time.Parse(util.TimeLayout, "08:00:00")

	selectedDateSlice := []time.Time{
		selectedDate,
	}

	params := GetAvailableTimeParams{
		PlaceID:      1,
		SelectedDate: selectedDate,
		StartTime:    startTime,
		BookedSlot:   10,
	}

	midnight := time.Date(params.SelectedDate.Year(), params.SelectedDate.Month(), params.SelectedDate.Day(), 0, 0, 0, 0, params.SelectedDate.Location())
	midnight = midnight.Add(time.Duration(1*24) * time.Hour)

	repoParams := GetBookingDataParams{
		PlaceID:   params.PlaceID,
		StartDate: params.SelectedDate,
		EndDate:   midnight,
		StartTime: params.StartTime,
	}

	var getBookingDataReturn []DataForCheckAvailableSchedule
	var getTimeSlotReturn []TimeSlot

	t.Run("get time slot data failed", func(t *testing.T) {
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", params.PlaceID, selectedDateSlice).Return(&getTimeSlotReturn, errors.Wrap(ErrInternalServerError, "test error"))

		output, err := mockService.GetAvailableTime(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestService_GetAvailableTimeGetPlaceCapacityFailedInternalServerError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockXenditService := new(MockXenditService)
	mockService := NewService(mockRepo, mockXenditService)

	selectedDate, _ := time.Parse(util.DateLayout, "2022-03-29")
	startTime, _ := time.Parse(util.TimeLayout, "08:00:00")

	selectedDateSlice := []time.Time{
		selectedDate,
	}

	params := GetAvailableTimeParams{
		PlaceID:      1,
		SelectedDate: selectedDate,
		StartTime:    startTime,
		BookedSlot:   10,
	}

	midnight := time.Date(params.SelectedDate.Year(), params.SelectedDate.Month(), params.SelectedDate.Day(), 0, 0, 0, 0, params.SelectedDate.Location())
	midnight = midnight.Add(time.Duration(1*24) * time.Hour)

	repoParams := GetBookingDataParams{
		PlaceID:   params.PlaceID,
		StartDate: params.SelectedDate,
		EndDate:   midnight,
		StartTime: params.StartTime,
	}

	var getBookingDataReturn []DataForCheckAvailableSchedule
	var getTimeSlotReturn []TimeSlot

	t.Run("get place capacity data failed", func(t *testing.T) {
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", params.PlaceID, selectedDateSlice).Return(&getTimeSlotReturn, nil)
		mockRepo.On("GetPlaceCapacity", params.PlaceID).Return(&PlaceOpenHourAndCapacity{}, errors.Wrap(ErrInternalServerError, "test error"))

		output, err := mockService.GetAvailableTime(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestService_GetAvailableDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "08:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingCapacityOne := 20

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingCapacityTwo := 20

		params := GetAvailableDateParams{
			PlaceID:    1,
			StartDate:  startDate,
			Interval:   3,
			BookedSlot: 10,
		}

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.StartDate,
			EndDate:   params.StartDate.Add(time.Duration(params.Interval*24) * time.Hour),
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      startDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      startDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		timeSlotOneStartTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotTwoStartTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotThreeStartTime, _ := time.Parse(util.TimeLayout, "10:00:00")

		getTimeSlotReturn := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotThreeStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       3,
			},
		}

		placeCapacity := 20

		checkedDateOne, _ := time.Parse(util.DateLayout, "2022-03-29")
		checkedDateTwo, _ := time.Parse(util.DateLayout, "2022-03-30")
		checkedDateThree, _ := time.Parse(util.DateLayout, "2022-03-31")
		checkedDateFour, _ := time.Parse(util.DateLayout, "2022-04-01")
		checkedDateSlice := []time.Time{
			checkedDateOne, checkedDateTwo, checkedDateThree, checkedDateFour,
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", repoParams.PlaceID, checkedDateSlice).Return(&getTimeSlotReturn, nil)
		mockRepo.On("GetPlaceCapacity", repoParams.PlaceID).Return(&PlaceOpenHourAndCapacity{OpenHour: timeSlotOneStartTime, Capacity: placeCapacity}, nil)

		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, output)
	})

	t.Run("success with default value of interval", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "08:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingCapacityOne := 20

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingCapacityTwo := 20

		params := GetAvailableDateParams{
			PlaceID:    1,
			StartDate:  startDate,
			Interval:   0,
			BookedSlot: 10,
		}

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.StartDate,
			EndDate:   params.StartDate.Add(time.Duration(7*24) * time.Hour),
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      startDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      startDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		timeSlotOneStartTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotTwoStartTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotThreeStartTime, _ := time.Parse(util.TimeLayout, "10:00:00")

		getTimeSlotReturn := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotTwoStartTime,
				EndTime:   timeSlotThreeStartTime,
				Day:       2,
			},
		}

		placeCapacity := 20

		checkedDateOne, _ := time.Parse(util.DateLayout, "2022-03-29")
		checkedDateTwo, _ := time.Parse(util.DateLayout, "2022-03-30")
		checkedDateThree, _ := time.Parse(util.DateLayout, "2022-03-31")
		checkedDateFour, _ := time.Parse(util.DateLayout, "2022-04-01")
		checkedDateFive, _ := time.Parse(util.DateLayout, "2022-04-02")
		checkedDateSix, _ := time.Parse(util.DateLayout, "2022-04-03")
		checkedDateSeven, _ := time.Parse(util.DateLayout, "2022-04-04")
		checkedDateEight, _ := time.Parse(util.DateLayout, "2022-04-05")
		checkedDateSlice := []time.Time{
			checkedDateOne, checkedDateTwo, checkedDateThree, checkedDateFour, checkedDateFive, checkedDateSix, checkedDateSeven, checkedDateEight,
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", repoParams.PlaceID, checkedDateSlice).Return(&getTimeSlotReturn, nil)
		mockRepo.On("GetPlaceCapacity", repoParams.PlaceID).Return(&PlaceOpenHourAndCapacity{OpenHour: timeSlotOneStartTime, Capacity: placeCapacity}, nil)

		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, output)
	})

	t.Run("failed get booking data", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityOne := 10

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityTwo := 10

		params := GetAvailableDateParams{
			PlaceID:    1,
			StartDate:  startDate,
			Interval:   3,
			BookedSlot: 10,
		}

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.StartDate,
			EndDate:   params.StartDate.Add(time.Duration(params.Interval*24) * time.Hour),
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      startDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      startDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, errors.Wrap(ErrInternalServerError, "test error"))
		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
	})

	t.Run("failed get time slot data", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityOne := 10

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityTwo := 10

		params := GetAvailableDateParams{
			PlaceID:    1,
			StartDate:  startDate,
			Interval:   3,
			BookedSlot: 10,
		}

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.StartDate,
			EndDate:   params.StartDate.Add(time.Duration(params.Interval*24) * time.Hour),
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      startDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      startDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		timeSlotOneStartTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotTwoStartTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotThreeStartTime, _ := time.Parse(util.TimeLayout, "10:00:00")
		timeSlotFourStartTime, _ := time.Parse(util.TimeLayout, "11:00:00")

		getTimeSlotReturn := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotTwoStartTime,
				EndTime:   timeSlotThreeStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotThreeStartTime,
				EndTime:   timeSlotFourStartTime,
				Day:       2,
			},
		}

		checkedDateOne, _ := time.Parse(util.DateLayout, "2022-03-29")
		checkedDateTwo, _ := time.Parse(util.DateLayout, "2022-03-30")
		checkedDateThree, _ := time.Parse(util.DateLayout, "2022-03-31")
		checkedDateFour, _ := time.Parse(util.DateLayout, "2022-04-01")
		checkedDateSlice := []time.Time{
			checkedDateOne, checkedDateTwo, checkedDateThree, checkedDateFour,
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", repoParams.PlaceID, checkedDateSlice).Return(&getTimeSlotReturn, errors.Wrap(ErrInternalServerError, "test error"))

		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
	})

	t.Run("failed get place capacity", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		bookingStartTimeOne, _ := time.Parse(util.TimeLayout, "09:00:00")
		bookingStartEndOne, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityOne := 10

		bookingStartTimeTwo, _ := time.Parse(util.TimeLayout, "10:00:00")
		bookingStartEndTwo, _ := time.Parse(util.TimeLayout, "11:00:00")
		bookingCapacityTwo := 10

		params := GetAvailableDateParams{
			PlaceID:    1,
			StartDate:  startDate,
			Interval:   3,
			BookedSlot: 10,
		}

		repoParams := GetBookingDataParams{
			PlaceID:   params.PlaceID,
			StartDate: params.StartDate,
			EndDate:   params.StartDate.Add(time.Duration(params.Interval*24) * time.Hour),
		}

		getBookingDataReturn := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      startDate,
				StartTime: bookingStartTimeOne,
				EndTime:   bookingStartEndOne,
				Capacity:  bookingCapacityOne,
			},
			{
				ID:        2,
				Date:      startDate,
				StartTime: bookingStartTimeTwo,
				EndTime:   bookingStartEndTwo,
				Capacity:  bookingCapacityTwo,
			},
		}

		timeSlotOneStartTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		timeSlotTwoStartTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		timeSlotThreeStartTime, _ := time.Parse(util.TimeLayout, "10:00:00")
		timeSlotFourStartTime, _ := time.Parse(util.TimeLayout, "11:00:00")

		getTimeSlotReturn := []TimeSlot{
			{
				ID:        1,
				StartTime: timeSlotOneStartTime,
				EndTime:   timeSlotTwoStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotTwoStartTime,
				EndTime:   timeSlotThreeStartTime,
				Day:       2,
			},
			{
				ID:        1,
				StartTime: timeSlotThreeStartTime,
				EndTime:   timeSlotFourStartTime,
				Day:       2,
			},
		}

		placeCapacity := 20

		checkedDateOne, _ := time.Parse(util.DateLayout, "2022-03-29")
		checkedDateTwo, _ := time.Parse(util.DateLayout, "2022-03-30")
		checkedDateThree, _ := time.Parse(util.DateLayout, "2022-03-31")
		checkedDateFour, _ := time.Parse(util.DateLayout, "2022-04-01")
		checkedDateSlice := []time.Time{
			checkedDateOne, checkedDateTwo, checkedDateThree, checkedDateFour,
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingDataReturn, nil)
		mockRepo.On("GetTimeSlotsData", repoParams.PlaceID, checkedDateSlice).Return(&getTimeSlotReturn, nil)
		mockRepo.On("GetPlaceCapacity", repoParams.PlaceID).Return(&PlaceOpenHourAndCapacity{OpenHour: timeSlotOneStartTime, Capacity: placeCapacity}, errors.Wrap(ErrInternalServerError, "test error"))

		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
	})

	t.Run("input validation error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)
		mockService := NewService(mockRepo, mockXenditService)

		startDate, _ := time.Parse(util.DateLayout, "2022-03-29")

		params := GetAvailableDateParams{
			PlaceID:    0,
			StartDate:  startDate,
			Interval:   0,
			BookedSlot: 0,
		}

		output, err := mockService.GetAvailableDate(params)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, output)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestService_CreateBooking(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             1,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 1,
			},
			{
				ID:      5,
				PlaceID: 1,
			},
		}

		bookingParams := CreateBookingParams{
			UserID:     input.UserID,
			PlaceID:    input.PlaceID,
			Date:       input.Date,
			StartTime:  input.StartTime,
			EndTime:    input.EndTime,
			Capacity:   input.Count,
			Status:     util.BookingMenungguKonfirmasi,
			TotalPrice: 0,
		}

		bookingItemParams := []CreateBookingItemsParams{
			{
				BookingID:  1,
				ItemID:     4,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  1,
				ItemID:     5,
				TotalPrice: 20000,
				Qty:        2,
			},
		}

		updateTotalPrice := UpdateTotalPriceParams{
			BookingID:  1,
			TotalPrice: 40000,
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 100,
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)
		mockRepo.On("CreateBooking", bookingParams).Return(&CreateBookingResponse{ID: 1}, nil)
		mockRepo.On("CreateBookingItems", bookingItemParams).Return(&CreateBookingItemsResponse{TotalPrice: 40000}, nil)
		mockRepo.On("UpdateTotalPrice", updateTotalPrice).Return(true, nil)

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		mockXenditService.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.BookingID)
	})

	t.Run("failed count < 0", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               -1,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		resp, err := service.CreateBooking(input)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})

	t.Run("failed item validation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
			{
				ID:      5,
				PlaceID: 2,
			},
		}

		itemsOutput := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
		}

		mockRepo.On("CheckedItem", items).Return(&itemsOutput, false, errors.Wrap(ErrInputValidationError, "test error"))

		resp, err := service.CreateBooking(input)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})

	t.Run("failed item check internal server error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
			{
				ID:      5,
				PlaceID: 2,
			},
		}

		itemsOutput := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
		}

		mockRepo.On("CheckedItem", items).Return(&itemsOutput, false, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed create booking internal server error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
			{
				ID:      5,
				PlaceID: 2,
			},
		}

		bookingParams := CreateBookingParams{
			UserID:     input.UserID,
			PlaceID:    input.PlaceID,
			Date:       input.Date,
			StartTime:  input.StartTime,
			EndTime:    input.EndTime,
			Capacity:   input.Count,
			Status:     util.BookingMenungguKonfirmasi,
			TotalPrice: 0,
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 100,
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)
		mockRepo.On("CreateBooking", bookingParams).Return(&CreateBookingResponse{ID: 1}, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed create booking item", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
			{
				ID:      5,
				PlaceID: 2,
			},
		}

		bookingParams := CreateBookingParams{
			UserID:     input.UserID,
			PlaceID:    input.PlaceID,
			Date:       input.Date,
			StartTime:  input.StartTime,
			EndTime:    input.EndTime,
			Capacity:   input.Count,
			Status:     util.BookingMenungguKonfirmasi,
			TotalPrice: 0,
		}

		bookingItemParams := []CreateBookingItemsParams{
			{
				BookingID:  1,
				ItemID:     4,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  1,
				ItemID:     5,
				TotalPrice: 20000,
				Qty:        2,
			},
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 100,
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)
		mockRepo.On("CreateBooking", bookingParams).Return(&CreateBookingResponse{ID: 1}, nil)
		mockRepo.On("CreateBookingItems", bookingItemParams).Return(&CreateBookingItemsResponse{TotalPrice: 40000}, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed create booking item internal server error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             2,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 2,
			},
			{
				ID:      5,
				PlaceID: 2,
			},
		}

		bookingParams := CreateBookingParams{
			UserID:     input.UserID,
			PlaceID:    input.PlaceID,
			Date:       input.Date,
			StartTime:  input.StartTime,
			EndTime:    input.EndTime,
			Capacity:   input.Count,
			Status:     util.BookingMenungguKonfirmasi,
			TotalPrice: 0,
		}

		bookingItemParams := []CreateBookingItemsParams{
			{
				BookingID:  1,
				ItemID:     4,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  1,
				ItemID:     5,
				TotalPrice: 20000,
				Qty:        2,
			},
		}

		updateTotalPrice := UpdateTotalPriceParams{
			BookingID:  1,
			TotalPrice: 40000,
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 100,
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)
		mockRepo.On("CreateBooking", bookingParams).Return(&CreateBookingResponse{ID: 1}, nil)
		mockRepo.On("CreateBookingItems", bookingItemParams).Return(&CreateBookingItemsResponse{TotalPrice: 40000}, nil)
		mockRepo.On("UpdateTotalPrice", updateTotalPrice).Return(false, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("success no item", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items:               []Item{},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             1,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		bookingParams := CreateBookingParams{
			UserID:     input.UserID,
			PlaceID:    input.PlaceID,
			Date:       input.Date,
			StartTime:  input.StartTime,
			EndTime:    input.EndTime,
			Capacity:   input.Count,
			Status:     util.BookingMenungguKonfirmasi,
			TotalPrice: 0,
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 100,
		}

		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)
		mockRepo.On("CreateBooking", bookingParams).Return(&CreateBookingResponse{ID: 1}, nil)

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		mockXenditService.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, resp.BookingID)
	})

	t.Run("failed when called get available time service", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             1,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 1,
			},
			{
				ID:      5,
				PlaceID: 1,
			},
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		mockXenditService.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed when called get available time service", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             1,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 1,
			},
			{
				ID:      5,
				PlaceID: 1,
			},
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)
		mockXenditService.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed when get available time result not match", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXenditService := new(MockXenditService)

		service := NewService(mockRepo, mockXenditService)

		date, _ := time.Parse(util.DateLayout, "2022-02-02")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		EndTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		input := CreateBookingServiceRequest{
			Items: []Item{
				{
					ID:    4,
					Name:  "Test Item 1",
					Price: 10000,
					Qty:   2,
				},
				{
					ID:    5,
					Name:  "Test Item 2",
					Price: 10000,
					Qty:   2,
				},
			},
			Date:                date,
			StartTime:           startTime,
			EndTime:             EndTime,
			Count:               10,
			PlaceID:             1,
			UserID:              1,
			CustomerName:        "Rafi Muhammad",
			CustomerPhoneNumber: "081291264758",
		}

		items := []CheckedItemParams{
			{
				ID:      4,
				PlaceID: 1,
			},
			{
				ID:      5,
				PlaceID: 1,
			},
		}

		midnight := time.Date(input.Date.Year(), input.Date.Month(), input.Date.Day(), 0, 0, 0, 0, input.Date.Location())
		midnight = midnight.Add(time.Duration(1*24) * time.Hour)

		repoParams := GetBookingDataParams{
			PlaceID:   input.PlaceID,
			StartDate: input.Date,
			EndDate:   midnight,
			StartTime: input.StartTime,
		}

		getBookingData := []DataForCheckAvailableSchedule{
			{
				ID:        1,
				Date:      input.Date,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Capacity:  input.Count,
			},
		}

		timeSlotsData := []TimeSlot{
			{
				ID:        1,
				StartTime: input.StartTime,
				EndTime:   input.EndTime,
				Day:       int(input.Date.Weekday()),
			},
		}

		placeIDAndCapacity := PlaceOpenHourAndCapacity{
			OpenHour: input.StartTime,
			Capacity: 0,
		}

		mockRepo.On("CheckedItem", items).Return(&items, true, nil)
		mockRepo.On("GetBookingData", repoParams).Return(&getBookingData, nil)
		mockRepo.On("GetTimeSlotsData", input.PlaceID, []time.Time{input.Date}).Return(&timeSlotsData, nil)
		mockRepo.On("GetPlaceCapacity", input.PlaceID).Return(&placeIDAndCapacity, nil)

		resp, err := service.CreateBooking(input)
		mockRepo.AssertExpectations(t)

		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestService_GetTimeSlots(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)
		service := NewService(repo, xenditService)

		date, _ := time.Parse(util.DateLayout, "2020-01-01")
		dateSlice := []time.Time{date}
		eightTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		nineTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		timeSlots := []TimeSlot{
			{
				ID:        1,
				StartTime: eightTime,
				EndTime:   nineTime,
				Day:       0,
			},
		}

		repo.On("GetTimeSlotsData", 1, dateSlice).Return(&timeSlots, nil)

		slots, err := service.GetTimeSlots(1, date)
		repo.AssertExpectations(t)
		assert.NotNil(t, slots)
		assert.Nil(t, err)
		assert.Equal(t, &timeSlots, slots)
	})

	t.Run("failed input validation error", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)
		service := NewService(repo, xenditService)

		date, _ := time.Parse(util.DateLayout, "2020-01-01")
		dateSlice := []time.Time{date}
		eightTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		nineTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		timeSlots := []TimeSlot{
			{
				ID:        1,
				StartTime: eightTime,
				EndTime:   nineTime,
				Day:       0,
			},
		}

		repo.On("GetTimeSlotsData", 1, dateSlice).Return(&timeSlots, nil)

		slots, err := service.GetTimeSlots(-1, date)
		assert.Nil(t, slots)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})

	t.Run("failed internal server error", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)
		service := NewService(repo, xenditService)

		date, _ := time.Parse(util.DateLayout, "2020-01-01")
		dateSlice := []time.Time{date}
		eightTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		nineTime, _ := time.Parse(util.TimeLayout, "09:00:00")

		timeSlots := []TimeSlot{
			{
				ID:        1,
				StartTime: eightTime,
				EndTime:   nineTime,
				Day:       0,
			},
		}

		repo.On("GetTimeSlotsData", 1, dateSlice).Return(&timeSlots, errors.Wrap(ErrInternalServerError, "test error"))

		slots, err := service.GetTimeSlots(1, date)
		repo.AssertExpectations(t)
		assert.Nil(t, slots)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestService_UpdateBookingStatusByXendit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)

		service := NewService(repo, xenditService)

		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "1",
			Status:     "PAID",
		}

		repo.On("UpdateBookingStatusByXenditID", "1", 2).
			Return(nil)

		err := service.UpdateBookingStatusByXendit(params)
		assert.Nil(t, err)
	})

	t.Run("failed from repo", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)

		service := NewService(repo, xenditService)

		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "1",
			Status:     "PAID",
		}

		repo.On("UpdateBookingStatusByXenditID", "1", 2).
			Return(ErrInternalServerError)

		err := service.UpdateBookingStatusByXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("success booking expired", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)

		service := NewService(repo, xenditService)

		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "1",
			Status:     "EXPIRED",
		}

		repo.On("UpdateBookingStatusByXenditID", "1", 4).
			Return(ErrInternalServerError)

		err := service.UpdateBookingStatusByXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("validation error unknown status", func(t *testing.T) {
		repo := new(MockRepository)
		xenditService := new(MockXenditService)

		service := NewService(repo, xenditService)

		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "1",
			Status:     "UNKNOWN",
		}

		err := service.UpdateBookingStatusByXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}
