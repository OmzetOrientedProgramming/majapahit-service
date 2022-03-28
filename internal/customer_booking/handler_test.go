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
	"github.com/pkg/errors"
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
		State:   1,
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listCustomerBooking := List{
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
			"bookings":   listCustomerBooking.CustomerBookings,
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

func TestHandler_GetListCustromerBookingWithPaginationPlaceIDError(t *testing.T) {
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
	ctx.SetParamValues("test")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect place id",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetListCustomerBookingWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemWithPaginationStateAndLimitAndPageAreNotInt(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "asd")
	q.Set("limit", "asd")
	q.Set("page", "asd")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/booking")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"state should be positive integer",
			"limit should be positive integer",
			"page should be positive integer",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetListCustomerBookingWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationWithStateLimitPageAreEmpty(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/booking")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListRequest{
		Limit:   0,
		Page:    0,
		Path:    "/api/v1/business-admin/1/booking",
		State:   0,
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listCustomerBooking := List{
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
			"bookings":   listCustomerBooking.CustomerBookings,
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

func TestHandler_GetListCustomerBookingWithPaginationLimitError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "0")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/booking")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/business-admin/1/booking",
		State:   0,
		PlaceID: 1,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listCustomerBooking List
	var pagination util.Pagination
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, errorFromService)

	// Tes
	assert.NoError(t, h.GetListCustomerBookingWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "0")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/booking")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/business-admin/1/booking",
		State:   0,
		PlaceID: 1,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listCustomerBooking List
	var pagination util.Pagination
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, internalServerError)

	// Tes
	assert.NoError(t, h.GetListCustomerBookingWithPagination(ctx))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}
