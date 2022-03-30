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
