package auth

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CheckPhoneNumber(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) VerifyOTP(phoneNumber, otp string) (bool, error) {
	args := m.Called(phoneNumber, otp)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) CreateCustomer(customer Customer) (*Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(*Customer), args.Error(1)
}

func (m *MockService) SendOTP(phoneNumber string) error {
	args := m.Called(phoneNumber)
	return args.Error(0)
}

func (m *MockService) Register(customer Customer) (*Customer, error) {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Customer), args.Error(1)
}

func (m *MockService) GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error) {
	args := m.Called(phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Customer), args.Error(1)
}

func TestHandler_CheckPhoneNumber(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	e := echo.New()

	t.Run("phone number is not registered", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "phone number is available",
			Data:    nil,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)
		mockService.On("SendOTP", "087748176534").Return(nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("incorrect query param", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"session must be register or login"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")

		req := httptest.NewRequest(http.MethodPost, "/?session=random", nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("incorrect request body", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")

		payload, _ := json.Marshal(map[string]int{
			"phone_number": 1,
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error on calling service function", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, ErrInternalServer)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("phone number already registered", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusConflict,
			Message: "phone number already registered",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed to send otp", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot send otp to phone number",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)
		mockService.On("SendOTP", "087748176534").Return(ErrInternalServer)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("phone number is not registered when log in", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusNotFound,
			Message: "phone number has not been registered",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("fail to send otp when log in", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot send otp to phone number",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		mockService.On("SendOTP", "087748176534").Return(ErrInternalServer)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("login success", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		mockService.On("SendOTP", "087748176534").Return(nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.CheckPhoneNumber(ctx))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_VerifyOTP(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	e := echo.New()

	t.Run("otp is valid", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    nil,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("VerifyOTP", "087748176534", "123456").Return(true, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
			"otp":          "123456",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.VerifyOTP(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("otp is not valid", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: "wrong otp code",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("VerifyOTP", "087748176534", "123456").Return(false, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
			"otp":          "123456",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.VerifyOTP(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("internal error on service layer to verify otp", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("VerifyOTP", "087748176534", "123456").Return(false, ErrInternalServer)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
			"otp":          "123456",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.VerifyOTP(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error binding request", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")

		payload, _ := json.Marshal(map[string]interface{}{
			"phone_number": "087748176534",
			"otp":          111111,
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.VerifyOTP(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("incorrect session query parameter", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"session must be register or login"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")

		req := httptest.NewRequest(http.MethodPost, "/?session=random", nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.VerifyOTP(ctx))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_Register(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      1,
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: "08123456789",
		}).Return(expectedCustomer, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "08123456789",
			"full_name":    "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()

		assert.NoError(t, mockHandler.Register(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error binding request to struct", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})

		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      1,
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: "08123456789",
		}).Return(expectedCustomer, nil)

		payload, _ := json.Marshal(map[string]interface{}{
			"phone_number": 123,
			"full_name":    "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.Register(ctx))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error on service layer", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService, "mockJWTSecret")
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: "08123456789",
		}).Return(nil, ErrInternalServer)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "08123456789",
			"full_name":    "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, mockHandler.Register(ctx))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
