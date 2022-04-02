package booking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
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

func (m *MockService) UpdateBookingStatus(bookingID int, newStatus int) error {
	args := m.Called(bookingID, newStatus)
	return args.Error(0)
}

func TestHandler_GetDetailSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID")
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
	c.SetPath("/api/v1/business-admin/booking/:bookingID")
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

func TestBooking_GetBookingDetailWithBookingIDBelowOne(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("0")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input
	bookingID := 0

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"placeID must be above 0"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}
	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var bookingDetail Detail
	mockService.On("GetDetail", bookingID).Return(&bookingDetail, errorFromService)

	// Test
	assert.NoError(t, h.GetDetail(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetDetailWithBookingIDString(t *testing.T) {
	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("satu")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"bookingID must be number",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetDetail(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_UpdateBookingStatusSuccess(t *testing.T) {
	// Setting up echo
	e := echo.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"status": 2,
	})

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID/confirmation")
	c.SetParamNames("bookingID")
	c.SetParamValues("1")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	bookingID := 1
	newStatus := 2

	// Expectation
	mockService.On("UpdateBookingStatus", bookingID, newStatus).Return(nil)

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "Success update status",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	if assert.NoError(t, h.UpdateBookingStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_UpdateBookingStatusWithBookingIDString(t *testing.T) {
	// Setting up echo
	e := echo.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"status": 2,
	})

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID/confirmation")
	c.SetParamNames("bookingID")
	c.SetParamValues("satu")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"bookingID must be number",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.UpdateBookingStatus(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_UpdateBookingStatusBindingError(t *testing.T) {
	// Setting up echo
	e := echo.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"status": "halo halo bandung",
	})

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID/confirmation")
	c.SetParamNames("bookingID")
	c.SetParamValues("1")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "cannot process request",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.UpdateBookingStatus(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_UpdateBookingStatusWithBookingIDBelowOne(t *testing.T) {
	// Setting up echo
	e := echo.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"status": 2,
	})

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID/confirmation")
	c.SetParamNames("bookingID")
	c.SetParamValues("0")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	bookingID := 0
	newStatus := 2

	// Expectation
	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"bookingID must be above 0"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	mockService.On("UpdateBookingStatus", bookingID, newStatus).Return(errorFromService)

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.UpdateBookingStatus(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))

}

func TestHandler_UpdateBookingStatusWithInternalServerError(t *testing.T) {
	// Setting up echo
	e := echo.New()

	payload, _ := json.Marshal(map[string]interface{}{
		"status": 2,
	})

	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/booking/:bookingID/confirmation")
	c.SetParamNames("bookingID")
	c.SetParamValues("10")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)
	// Define input and output
	bookingID := 10
	newStatus := 2

	errorFromService := errors.Wrap(ErrInternalServerError, "test error")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("UpdateBookingStatus", bookingID, newStatus).Return(errorFromService)

	// Tes
	assert.NoError(t, h.UpdateBookingStatus(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (m *MockService) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	args := m.Called(localID)
	myBookingsOngoing := args.Get(0).(*[]Booking)
	return myBookingsOngoing, args.Error(1)
}

func TestHandler_GetMyBookingsOngoingSuccess(t *testing.T) {
	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userData", &userData)
	c.SetPath("/api/v1/booking/ongoing")

	// Setting up service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setting up Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Setting up input and output
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

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    myBookingsOngoing,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetMyBookingsOngoing", userData.Users[0].LocalID).Return(&myBookingsOngoing, nil)

	// Test Fields
	if assert.NoError(t, h.GetMyBookingsOngoing(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}
func TestService_GetMyBookingsOngoingWithEmptyLocalID(t *testing.T) {
	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userData", &userData)
	c.SetPath("/api/v1/booking/ongoing")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input
	localID := ""

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"localID cannot be empty"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}
	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var myBookingsOngoing []Booking
	mockService.On("GetMyBookingsOngoing", localID).Return(&myBookingsOngoing, errorFromService)

	// Test
	assert.NoError(t, h.GetMyBookingsOngoing(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func (m *MockService) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error) {
	args := m.Called(params)
	myBookingsPrevious := args.Get(0).(*List)
	pagination := args.Get(1).(util.Pagination)
	return myBookingsPrevious, &pagination, args.Error(2)
}

func TestHandler_GetMyBookingsPreviousWithPaginationWithParams(t *testing.T) {
	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userData", &userData)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/v1/booking/previous",
	}

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

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":   myBookingsPrevious.Bookings,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetMyBookingsPreviousWithPagination", params).Return(&myBookingsPrevious, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetMyBookingsPreviousWithPagination(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetMyBookingsPreviousWithPaginationWithValidationErrorLimitPageNotInt(t *testing.T) {
	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	// Setup echo
	e := echo.New()
	q := make(url.Values)
	q.Set("limit", "testerror")
	q.Set("page", "testerror")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userData", &userData)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"limit should be positive integer",
			"page should be positive integer",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetMyBookingsPreviousWithPagination(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetMyBookingsPreviousWithPaginationWithoutParams(t *testing.T) {
	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userData", &userData)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := BookingsListRequest{
		Limit: 0,
		Page:  0,
		Path:  "/api/v1/booking/previous",
	}

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

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/booking/previous?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":   myBookingsPrevious.Bookings,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetMyBookingsPreviousWithPagination", params).Return(&myBookingsPrevious, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetMyBookingsPreviousWithPagination(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}
