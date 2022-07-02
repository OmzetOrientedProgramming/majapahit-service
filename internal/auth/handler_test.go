package auth

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CheckPhoneNumber(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) VerifyOTP(sessionInfo, otp string) (*VerifyOTPResult, error) {
	args := m.Called(sessionInfo, otp)
	return args.Get(0).(*VerifyOTPResult), args.Error(1)
}

func (m *MockService) CreateCustomer(customer Customer) (*Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(*Customer), args.Error(1)
}

func (m *MockService) SendOTP(phoneNumber, recaptchaToken string) (string, error) {
	args := m.Called(phoneNumber, recaptchaToken)
	return args.String(0), args.Error(1)
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
			Message: "success",
			Data: CheckPhoneNumberResponse{
				SessionInfo: "test session token",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)
		mockService.On("SendOTP", "087748176534", "testToken").Return("test session token", nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "testToken",
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

	t.Run("validation error from sendOTP service register", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"test input validation error",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)
		mockService.On("SendOTP", "087748176534", "testToken").Return("", errors.Wrap(ErrInputValidation, "test input validation error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "testToken",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("validation error from sendOTP service login", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"test input validation error",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		mockService.On("SendOTP", "087748176534", "testToken").Return("", errors.Wrap(ErrInputValidation, "test input validation error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "testToken",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("incorrect query param", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"session must be register or login"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/?session=random", nil)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("incorrect request body", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		payload, _ := json.Marshal(map[string]int{
			"phone_number": 1,
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error on calling service function", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, errors.Wrap(ErrInternalServer, "test error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("phone number already registered", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"phone number already registered"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed to send otp", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)
		mockService.On("SendOTP", "087748176534", "test token").Return("", errors.Wrap(ErrInternalServer, "test error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "test token",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("phone number is not registered when log in", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"phone number has not been registered"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(false, nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "087748176534",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("fail to send otp when log in", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		mockService.On("SendOTP", "087748176534", "test token").Return("", errors.Wrap(ErrInternalServer, "test error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "test token",
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=login", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.CheckPhoneNumber(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("login success", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    CheckPhoneNumberResponse{SessionInfo: "test sessionToken"},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("CheckPhoneNumber", "087748176534").Return(true, nil)
		mockService.On("SendOTP", "087748176534", "test token").Return("test sessionToken", nil)

		payload, _ := json.Marshal(map[string]string{
			"phone_number":    "087748176534",
			"recaptcha_token": "test token",
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

	resp := VerifyOTPResult{
		AccessToken:  "test access token",
		RefreshToken: "test refresh token",
		ExpiresIn:    "300",
		LocalID:      "test local id",
		IsNewUser:    false,
		PhoneNumber:  "test phone number",
	}

	t.Run("otp is valid", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    resp,
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("VerifyOTP", "test session info", "test otp").Return(&resp, nil)

		payload, _ := json.Marshal(map[string]string{
			"session_info": "test session info",
			"otp":          "test otp",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
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
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"test input validation error",
			},
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("VerifyOTP", "test session info", "test otp").Return(&resp, errors.Wrap(ErrInputValidation, "test input validation error"))

		payload, _ := json.Marshal(map[string]string{
			"session_info": "test session info",
			"otp":          "test otp",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.VerifyOTP(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("internal error on service layer to verify otp", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("VerifyOTP", "test session info", "test otp").Return(&resp, errors.Wrap(ErrInternalServer, "test error"))

		payload, _ := json.Marshal(map[string]string{
			"session_info": "test session info",
			"otp":          "test otp",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.VerifyOTP(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error binding request", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		payload, _ := json.Marshal(map[string]interface{}{
			"phone_number": "087748176534",
			"otp":          111111,
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		util.ErrorHandler(mockHandler.VerifyOTP(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_Register(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
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

	t.Run("error binding request to struct", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      "customer",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: userData.Users[0].PhoneNumber,
			LocalID:     userData.Users[0].LocalID,
		}).Return(expectedCustomer, nil)

		payload, _ := json.Marshal(map[string]interface{}{
			"full_name": 123412,
		})
		req := httptest.NewRequest(http.MethodPost, "/?session=register", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.Register(c), c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("success", func(t *testing.T) {
		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      "customer",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: userData.Users[0].PhoneNumber,
			LocalID:     userData.Users[0].LocalID,
		}).Return(expectedCustomer, nil)

		payload, _ := json.Marshal(map[string]string{
			"full_name": "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Bearer token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)

		assert.NoError(t, mockHandler.Register(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("error input validation", func(t *testing.T) {
		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      "customer",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: userData.Users[0].PhoneNumber,
			LocalID:     userData.Users[0].LocalID,
		}).Return(expectedCustomer, errors.Wrap(ErrInputValidation, "test input validation error"))

		payload, _ := json.Marshal(map[string]string{
			"full_name": "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Bearer token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.Register(c), c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
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

		expectedCustomer := &Customer{
			ID:          1,
			Name:        "customer name",
			PhoneNumber: "08123456789",
			Status:      "customer",
		}

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: userDataFailed.Users[0].PhoneNumber,
			LocalID:     userDataFailed.Users[0].LocalID,
		}).Return(expectedCustomer, errors.Wrap(ErrInputValidation, "test input validation error"))

		payload, _ := json.Marshal(map[string]string{
			"full_name": "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Bearer token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userDataFailed)

		util.ErrorHandler(mockHandler.Register(c), c)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("error on service layer", func(t *testing.T) {
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)
		mockService.On("Register", Customer{
			Name:        "customer name",
			PhoneNumber: userData.Users[0].PhoneNumber,
			LocalID:     userData.Users[0].LocalID,
		}).Return(nil, errors.Wrap(ErrInternalServer, "test error"))

		payload, _ := json.Marshal(map[string]string{
			"phone_number": "08123456789",
			"full_name":    "customer name",
		})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.Register(c), c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
