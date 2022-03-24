package customerbooking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetListCustomerBookingWithPagination(params ListRequest) (*List, *util.Pagination, error) {
	args := m.Called(params)
	listCustomerBooking := args.Get(0).(*List)
	pagination := args.Get(1).(util.Pagination)
	return listCustomerBooking, &pagination, args.Error(2)
}

func TestHandler_GetListCustomerBookingWithPaginationSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "1")
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/booking")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/1/booking",
		State: 	1,
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listCustomerBooking := List{
		CustomerBookings: []CustomerBooking{
			{
				ID:          1,
				CustomerName: "test name 1",
				Capacity:     10,
				Date: "test date 1",
				StartTime: "test start time 1",
				EndTime: "test end time 1",
			},
			{
				ID:          2,
				CustomerName: "test name 2",
				Capacity:     10,
				Date: "test date 2",
				StartTime: "test start time 2",
				EndTime: "test end time 2",
			},
		},
		TotalCount: 10,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/1/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/1/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/1/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/1/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":      listCustomerBooking.CustomerBookings,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListCustomerBookingWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}