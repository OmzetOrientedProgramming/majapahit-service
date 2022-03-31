package booking

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetDetail(bookingID int) (*Detail, error) {
	args := m.Called(bookingID)
	bookingDetail := args.Get(0).(*Detail)
	return bookingDetail, args.Error(1)
}

func TestHandler_GetDetailSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	bookingID := 1

	createdAtRow := time.Date(2021, time.Month(10), 26, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	bookingDetail := Detail{
		ID:               1,
		Date:             "27 Oktober 2021",
		StartTime:        "19:00",
		EndTime:          "20:00",
		Capacity:         10,
		Status:           1,
		TotalPrice:       500000.0,
		TotalPriceItem:   400000.0,
		TotalPriceTicket: 100000.0,
		CreatedAt:        createdAtRow,
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

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    bookingDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetDetail", bookingID).Return(&bookingDetail, nil)

	// Test Fields
	if assert.NoError(t, h.GetDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetDetailWithInternalServerError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("10")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input and output
	bookingID := 10

	errorFromService := errors.Wrap(ErrInternalServerError, "test error")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var bookingDetail Detail
	mockService.On("GetDetail", bookingID).Return(&bookingDetail, errorFromService)

	// Tes
	assert.NoError(t, h.GetDetail(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}
