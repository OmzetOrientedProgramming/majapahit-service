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
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, *util.Pagination, error) {
	args := m.Called(params)
	listCustomerBooking := args.Get(0).(*ListBooking)
	pagination := args.Get(1).(util.Pagination)
	return listCustomerBooking, &pagination, args.Error(2)
}

func (m *MockService) GetAvailableTime(params GetAvailableTimeParams) (*[]AvailableTimeResponse, error) {
	args := m.Called(params)
	return args.Get(0).(*[]AvailableTimeResponse), args.Error(1)
}

func (m *MockService) GetAvailableDate(params GetAvailableDateParams) (*[]AvailableDateResponse, error) {
	args := m.Called(params)
	return args.Get(0).(*[]AvailableDateResponse), args.Error(1)
}

func (m *MockService) CreateBooking(params CreateBookingServiceRequest) (*CreateBookingServiceResponse, error) {
	args := m.Called(params)
	return args.Get(0).(*CreateBookingServiceResponse), args.Error(1)
}

func (m *MockService) GetTimeSlots(placeID int, selectedDate time.Time) (*[]TimeSlot, error) {
	args := m.Called(placeID, selectedDate)
	return args.Get(0).(*[]TimeSlot), args.Error(1)
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
func (m *MockService) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	args := m.Called(localID)
	myBookingsOngoing := args.Get(0).(*[]Booking)
	return myBookingsOngoing, args.Error(1)
}

func (m *MockService) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error) {
	args := m.Called(localID, params)
	myBookingsPrevious := args.Get(0).(*List)
	pagination := args.Get(1).(util.Pagination)
	return myBookingsPrevious, &pagination, args.Error(2)
}

func (m *MockService) UpdateBookingStatusByXendit(callback XenditInvoicesCallback) error {
	args := m.Called(callback)
	return args.Error(0)
}

func TestHandler_GetListCustomerBookingWithPaginationSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "1")
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/v1/business-admin/booking",
		State:  1,
		UserID: userModel.ID,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listCustomerBooking := ListBooking{
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
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
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

func TestHandler_GetListItemWithPaginationStateAndLimitAndPageAreNotInt(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "asd")
	q.Set("limit", "asd")
	q.Set("page", "asd")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

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
	util.ErrorHandler(h.GetListCustomerBookingWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationWithStateLimitPageAreEmpty(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/v1/business-admin/booking",
		State:  0,
		UserID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listCustomerBooking := ListBooking{
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
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/booking?state=1&limit=10&page=1", os.Getenv("BASE_URL")),
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

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "0")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/booking",
		State:  0,
		UserID: 1,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listCustomerBooking ListBooking
	var pagination util.Pagination
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, errorFromService)

	// Tes
	util.ErrorHandler(h.GetListCustomerBookingWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "0")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/booking",
		State:  0,
		UserID: 1,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listCustomerBooking ListBooking
	var pagination util.Pagination
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, internalServerError)

	// Tes
	util.ErrorHandler(h.GetListCustomerBookingWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
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

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("state", "0")
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/booking?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/booking",
		State:  0,
		UserID: 1,
	}

	// Excpectation
	var listCustomerBooking ListBooking
	var pagination util.Pagination
	mockService.On("GetListCustomerBookingWithPagination", params).Return(&listCustomerBooking, pagination, nil)

	// Tes
	util.ErrorHandler(h.GetListCustomerBookingWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetAvailableTime(t *testing.T) {
	e := echo.New()

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

	placeID := 1
	dateString := "2022-03-29"
	date, _ := time.Parse(util.DateLayout, dateString)
	checkInString := "08:00:00"
	checkIn, _ := time.Parse(util.TimeLayout, checkInString)
	count := 10

	params := GetAvailableTimeParams{
		PlaceID:      placeID,
		SelectedDate: date,
		StartTime:    checkIn,
		BookedSlot:   count,
	}

	returnedData := []AvailableTimeResponse{
		{
			Time:  "09:00:00",
			Total: count,
		},
	}

	t.Run("success", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    returnedData,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableTime", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("check_in", checkInString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		assert.NoError(t, mockHandler.GetAvailableTime(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error forbidden", func(t *testing.T) {
		userDataFailed := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
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
		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableTime", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("check_in", checkInString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userDataFailed)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableTime(ctx), ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("input validation error", func(t *testing.T) {
		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableTime", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "testWrong")
		q.Set("date", "20:20:20")
		q.Set("check_in", "2020-02-02")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("test")

		util.ErrorHandler(mockHandler.GetAvailableTime(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHandler_GetAvailableTimeInternalServerError(t *testing.T) {
	e := echo.New()

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

	placeID := 1
	dateString := "2022-03-29"
	date, _ := time.Parse(util.DateLayout, dateString)
	checkInString := "08:00:00"
	checkIn, _ := time.Parse(util.TimeLayout, checkInString)
	count := 10

	params := GetAvailableTimeParams{
		PlaceID:      placeID,
		SelectedDate: date,
		StartTime:    checkIn,
		BookedSlot:   count,
	}

	var returnedData []AvailableTimeResponse

	t.Run("internal server error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableTime", params).Return(&returnedData, errors.Wrap(ErrInternalServerError, "test error"))

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("check_in", checkInString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableTime(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("input validation error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"test error",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableTime", params).Return(&returnedData, errors.Wrap(ErrInputValidationError, "test error"))

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("check_in", checkInString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableTime(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_GetAvailableDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()

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

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		interval := "3"
		count := 10

		params := GetAvailableDateParams{
			PlaceID:    placeID,
			StartDate:  date,
			Interval:   3,
			BookedSlot: count,
		}

		returnedData := []AvailableDateResponse{
			{
				Date:   "09:00:00",
				Status: "available",
			},
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    returnedData,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableDate", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("interval", interval)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		assert.NoError(t, mockHandler.GetAvailableDate(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("internal server error from service", func(t *testing.T) {
		e := echo.New()

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

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		interval := "3"
		count := 10

		params := GetAvailableDateParams{
			PlaceID:    placeID,
			StartDate:  date,
			Interval:   3,
			BookedSlot: count,
		}

		returnedData := []AvailableDateResponse{
			{
				Date:   "09:00:00",
				Status: "available",
			},
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableDate", params).Return(&returnedData, errors.Wrap(ErrInternalServerError, "test error"))

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("interval", interval)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableDate(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error forbidden", func(t *testing.T) {
		e := echo.New()
		userDataFailed := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
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

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		interval := "3"
		count := 10

		params := GetAvailableDateParams{
			PlaceID:    placeID,
			StartDate:  date,
			Interval:   3,
			BookedSlot: count,
		}

		returnedData := []AvailableDateResponse{
			{
				Date:   "09:00:00",
				Status: "available",
			},
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableDate", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("interval", interval)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userDataFailed)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableDate(ctx), ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("input validation error", func(t *testing.T) {
		e := echo.New()
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

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		count := 10

		params := GetAvailableDateParams{
			PlaceID:    placeID,
			StartDate:  date,
			Interval:   3,
			BookedSlot: count,
		}

		returnedData := []AvailableDateResponse{
			{
				Date:   "09:00:00",
				Status: "available",
			},
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableDate", params).Return(&returnedData, nil)

		q := make(url.Values)
		q.Set("count", "testWrong")
		q.Set("date", "20:20:20")
		q.Set("interval", "2020-02-02")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("test")

		util.ErrorHandler(mockHandler.GetAvailableDate(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("input validation error from service", func(t *testing.T) {
		e := echo.New()

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

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		interval := "3"
		count := 10

		params := GetAvailableDateParams{
			PlaceID:    placeID,
			StartDate:  date,
			Interval:   3,
			BookedSlot: count,
		}

		returnedData := []AvailableDateResponse{
			{
				Date:   "09:00:00",
				Status: "available",
			},
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"test error",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetAvailableDate", params).Return(&returnedData, errors.Wrap(ErrInputValidationError, "test error"))

		q := make(url.Values)
		q.Set("count", "10")
		q.Set("date", dateString)
		q.Set("interval", interval)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetAvailableDate(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_CreateBooking(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      dateString,
			StartTime: startTimeString,
			EndTime:   endTimeString,
			Count:     count,
		}

		payload, _ := json.Marshal(params)

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    returnedData,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		assert.NoError(t, mockHandler.CreateBooking(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("forbidden", func(t *testing.T) {
		e := echo.New()

		userDataFailed := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      dateString,
			StartTime: startTimeString,
			EndTime:   endTimeString,
			Count:     count,
		}

		payload, _ := json.Marshal(params)

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userDataFailed)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.CreateBooking(ctx), ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("failed binding", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      dateString,
			StartTime: startTimeString,
			EndTime:   endTimeString,
			Count:     count,
		}

		payload, _ := json.Marshal(map[string]interface{}{
			"start_time": 1,
		})

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.CreateBooking(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("input validation error from request", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := 1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      "2022",
			StartTime: "08",
			EndTime:   "09",
			Count:     count,
		}

		payload, _ := json.Marshal(params)

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("testFailed")

		util.ErrorHandler(mockHandler.CreateBooking(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("place id < 0", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := -1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      dateString,
			StartTime: startTimeString,
			EndTime:   endTimeString,
			Count:     count,
		}

		payload, _ := json.Marshal(params)

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, errors.Wrap(ErrInputValidationError, "test error"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("-1")

		util.ErrorHandler(mockHandler.CreateBooking(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("internal server error from service", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		placeID := -1
		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
		startTimeString := "08:00:00"
		startTime, _ := time.Parse(util.TimeLayout, startTimeString)
		endTimeString := "09:00:00"
		endTime, _ := time.Parse(util.TimeLayout, endTimeString)
		count := 10

		params := CreateBookingRequestBody{
			Items: []Item{
				{
					ID:    1,
					Name:  "test item",
					Price: 10000,
					Qty:   1,
				},
			},
			Date:      dateString,
			StartTime: startTimeString,
			EndTime:   endTimeString,
			Count:     count,
		}

		payload, _ := json.Marshal(params)

		serviceRequest := CreateBookingServiceRequest{
			Items:               params.Items,
			Date:                date,
			StartTime:           startTime,
			EndTime:             endTime,
			Count:               params.Count,
			PlaceID:             placeID,
			UserID:              userFromDatabase.ID,
			CustomerName:        userFromDatabase.Name,
			CustomerPhoneNumber: userFromDatabase.PhoneNumber,
		}

		returnedData := CreateBookingServiceResponse{
			"10",
			1,
			"test.com",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("CreateBooking", serviceRequest).Return(&returnedData, errors.Wrap(ErrInternalServerError, "test error"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("-1")

		util.ErrorHandler(mockHandler.CreateBooking(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestHandler_GetTimeSlot(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
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

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetTimeSlots", 1, date).Return(&timeSlots, nil)

		q := make(url.Values)
		q.Set("date", dateString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		assert.NoError(t, mockHandler.GetTimeSlots(ctx))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("failed forbidden", func(t *testing.T) {
		e := echo.New()

		userDataFailed := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
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

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetTimeSlots", 1, date).Return(&timeSlots, nil)

		q := make(url.Values)
		q.Set("date", dateString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userDataFailed)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetTimeSlots(ctx), ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("input validation error from request", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
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

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetTimeSlots", 1, date).Return(&timeSlots, nil)

		q := make(url.Values)
		q.Set("date", "08")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("test")

		util.ErrorHandler(mockHandler.GetTimeSlots(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("input validation error from service", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
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

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetTimeSlots", 1, date).Return(&timeSlots, errors.Wrap(ErrInputValidationError, "test error"))

		q := make(url.Values)
		q.Set("date", dateString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetTimeSlots(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("internal server error from service", func(t *testing.T) {
		e := echo.New()

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

		userFromDatabase := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		dateString := "2022-03-29"
		date, _ := time.Parse(util.DateLayout, dateString)
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

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("GetTimeSlots", 1, date).Return(&timeSlots, errors.Wrap(ErrInternalServerError, "test error"))

		q := make(url.Values)
		q.Set("date", dateString)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromFirebase", &userData)
		ctx.Set("userFromDatabase", &userFromDatabase)
		ctx.SetPath("/booking/:placeID")
		ctx.SetParamNames("placeID")
		ctx.SetParamValues("1")

		util.ErrorHandler(mockHandler.GetTimeSlots(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
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
		TotalPrice:       415000.0,
		TotalPriceItem:   400000.0,
		TotalPriceTicket: 15000.0,
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
	util.ErrorHandler(h.GetDetail(c), c)
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
	util.ErrorHandler(h.GetDetail(c), c)
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
	util.ErrorHandler(h.GetDetail(c), c)
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
	util.ErrorHandler(h.UpdateBookingStatus(c), c)
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
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	util.ErrorHandler(h.UpdateBookingStatus(c), c)
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
	util.ErrorHandler(h.UpdateBookingStatus(c), c)
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
	util.ErrorHandler(h.UpdateBookingStatus(c), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
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
	c.Set("userFromFirebase", &userData)
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
	c.Set("userFromFirebase", &userData)
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
	util.ErrorHandler(h.GetMyBookingsOngoing(c), c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
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
	c.Set("userFromFirebase", &userData)

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

	localID := ""

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
	mockService.On("GetMyBookingsPreviousWithPagination", localID, params).Return(&myBookingsPrevious, pagination, nil)

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
	c.Set("userFromFirebase", &userData)

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
	util.ErrorHandler(h.GetMyBookingsPreviousWithPagination(c), c)
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
	c.Set("userFromFirebase", &userData)

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

	localID := ""

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
	mockService.On("GetMyBookingsPreviousWithPagination", localID, params).Return(&myBookingsPrevious, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetMyBookingsPreviousWithPagination(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_XenditInvoicesCallback(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "10",
			Status:     "PAID",
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusCreated,
			Message: "success",
		})

		payload, _ := json.Marshal(params)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UpdateBookingStatusByXendit", params).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.XenditInvoicesCallback(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed binding body", func(t *testing.T) {
		params := map[string]interface{}{
			"id":          1,
			"external_id": 110,
			"status":      0,
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		payload, _ := json.Marshal(params)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.XenditInvoicesCallback(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed input validation error from service", func(t *testing.T) {
		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "10",
			Status:     "PAID",
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"test error"},
		})

		payload, _ := json.Marshal(params)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UpdateBookingStatusByXendit", params).Return(errors.Wrap(ErrInputValidationError, "test error"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.XenditInvoicesCallback(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed internal server error from service", func(t *testing.T) {
		params := XenditInvoicesCallback{
			ID:         "1",
			ExternalID: "10",
			Status:     "PAID",
		}

		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		payload, _ := json.Marshal(params)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UpdateBookingStatusByXendit", params).Return(errors.Wrap(ErrInternalServerError, "test error"))

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.XenditInvoicesCallback(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
