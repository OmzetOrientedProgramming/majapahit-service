package businessadminauth

import (
	"bytes"
	"encoding/json"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Login(email, password, recaptchaToken string) (string, string, error) {
	args := m.Called(email, password, recaptchaToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockService) RegisterBusinessAdmin(request RegisterBusinessAdminRequest) (*LoginCredential, error) {
	args := m.Called(request)
	loginCredential := args.Get(0).(*LoginCredential)
	return loginCredential, args.Error(1)
}

func TestHandler_RegisterBusinessAdmin(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")

	mockRequest := RegisterBusinessAdminRequest{
		AdminPhoneNumber:        "089782828888",
		AdminEmail:              "sebuahemail@gmail.com",
		AdminName:               "Rafi Muhammad",
		AdminBankAccount:        "008-112492374950",
		AdminBankAccountName:    "RAFI MUHAMMAD",
		PlaceName:               "Kopi Kenangan",
		PlaceAddress:            "Jalan Raya Pasar Minggu",
		PlaceDescription:        "Kopi Kenangan menyediakan berbagai macam kopi sesuai pesanan Anda.",
		PlaceCapacity:           20,
		PlaceInterval:           30,
		PlaceImage:              "https://drive.google.com/file/d/.../view?usp=sharing",
		PlaceOpenHour:           "08:00",
		PlaceCloseHour:          "20:00",
		PlaceMinIntervalBooking: 1,
		PlaceMaxIntervalBooking: 3,
		PlaceMinSlotBooking:     1,
		PlaceMaxSlotBooking:     5,
		PlaceLat:                100.0,
		PlaceLong:               2.0002638,
	}

	loginCredentialExpected := &LoginCredential{
		PlaceName: "Kopi Kenangan",
		Email:     "sebuahemail@gmail.com",
		Password:  "12345678",
	}

	mockService := new(MockService)
	mockHandler := NewHandler(mockService)
	mockService.On("RegisterBusinessAdmin", mockRequest).Return(loginCredentialExpected, nil)

	request, _ := json.Marshal(mockRequest)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/auth/business-admin/register")

	assert.NoError(t, mockHandler.RegisterBusinessAdmin(e.NewContext(req, rec)))
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestHandler_Login(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")

	mockService := new(MockService)
	mockHandler := NewHandler(mockService)
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "Login berhasil",
			Data: LoginResponse{
				AccessToken:  "test access token",
				RefreshToken: "test refresh token",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.
			On("Login", "test@gmail.com", "testpass", "test captcha response").
			Return("test access token", "test refresh token", nil)

		payload, _ := json.Marshal(LoginRequest{
			Email:           "test@gmail.com",
			Password:        "testpass",
			CaptchaResponse: "test captcha response",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed to decode request", func(t *testing.T) {
		request, _ := json.Marshal(&map[string]interface{}{
			"email": 123,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(request))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("input validation error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "Kredensial yang anda berikan tidak valid",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.
			On("Login", "test@gmail.com", "testpass", "test captcha response").
			Return("test access token", "test refresh token", ErrInputValidationError)

		payload, _ := json.Marshal(LoginRequest{
			Email:           "test@gmail.com",
			Password:        "testpass",
			CaptchaResponse: "test captcha response",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("unauthorized error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Kredensial yang anda berikan salah",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.
			On("Login", "test@gmail.com", "testpass", "test captcha response").
			Return("test access token", "test refresh token", ErrUnauthorized)

		payload, _ := json.Marshal(LoginRequest{
			Email:           "test@gmail.com",
			Password:        "testpass",
			CaptchaResponse: "test captcha response",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("not found error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusNotFound,
			Message: "Akun tidak ditemukan",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.
			On("Login", "test@gmail.com", "testpass", "test captcha response").
			Return("test access token", "test refresh token", ErrNotFound)

		payload, _ := json.Marshal(LoginRequest{
			Email:           "test@gmail.com",
			Password:        "testpass",
			CaptchaResponse: "test captcha response",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("internal server error from service", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Terjadi kesalahan dalam memproses permintaan anda",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.
			On("Login", "test@gmail.com", "testpass", "test captcha response").
			Return("test access token", "test refresh token", ErrInternalServerError)

		payload, _ := json.Marshal(LoginRequest{
			Email:           "test@gmail.com",
			Password:        "testpass",
			CaptchaResponse: "test captcha response",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/business-admin/login", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/auth/business-admin/login")

		assert.NoError(t, mockHandler.Login(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
