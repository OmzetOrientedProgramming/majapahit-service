package booking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

type MockService struct {
	mock.Mock
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

func (m *MockService) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	args := m.Called(localID)
	myBookingsOngoing := args.Get(0).(*[]Booking)
	return myBookingsOngoing, args.Error(1)
}

func (m *MockService) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error) {
	args := m.Called(params)
	myBookingsPrevious := args.Get(0).(*List)
	pagination := args.Get(1).(util.Pagination)
	return myBookingsPrevious, &pagination, args.Error(2)
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

		assert.NoError(t, mockHandler.GetAvailableTime(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableTime(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableTime(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableTime(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableDate(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableDate(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableDate(ctx))
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

		assert.NoError(t, mockHandler.GetAvailableDate(ctx))
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

		assert.NoError(t, mockHandler.CreateBooking(ctx))
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

		assert.NoError(t, mockHandler.CreateBooking(ctx))
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

		assert.NoError(t, mockHandler.CreateBooking(ctx))
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

		assert.NoError(t, mockHandler.CreateBooking(ctx))
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

		assert.NoError(t, mockHandler.CreateBooking(ctx))
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

		assert.NoError(t, mockHandler.GetTimeSlots(ctx))
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

		assert.NoError(t, mockHandler.GetTimeSlots(ctx))
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

		assert.NoError(t, mockHandler.GetTimeSlots(ctx))
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

		assert.NoError(t, mockHandler.GetTimeSlots(ctx))
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

	// Expectation
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
	assert.NoError(t, h.GetMyBookingsOngoing(c))
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
	c.Set("userFromFirebase", &userData)

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
