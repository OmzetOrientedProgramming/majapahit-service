package customer

import (
	"bytes"
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
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) PutEditCustomer(body EditCustomerRequest) error {
	args := m.Called(body)
	return args.Error(0)
}

func TestHandler_PutEditCustomer(t *testing.T) {
	t.Run("Success PutEditCustomer", func(t *testing.T) {
		// Setup echo
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

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2001-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		body := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusOK,
			Message: "Successfully Edited Profile!",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("PutEditCustomer", body).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		assert.NoError(t, mockHandler.PutEditCustomer(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Error because of forbidden error", func(t *testing.T) {
		// Setup echo
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

		userModelFailed := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModelFailed.ID
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2001-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		body := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(body)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("PutEditCustomer", body).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Bearer token")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModelFailed)
		ctx.Set("userFromFirebase", &userDataFailed)

		util.ErrorHandler(mockHandler.PutEditCustomer(ctx), ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Binding error", func(t *testing.T) {
		// Setup echo
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

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		nameFailed := 123
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2001-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		type MismatchNameRequest struct {
			ID                int
			Name              int    `json:"name"`
			ProfilePicture    string `json:"profile_picture"`
			DateOfBirth       time.Time
			DateOfBirthString string `json:"date_of_birth"`
			Gender            int    `json:"gender"`
		}

		bodyFailed := MismatchNameRequest{
			ID:                userID,
			Name:              nameFailed,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(bodyFailed)

		expectedResponse := util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.PutEditCustomer(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Format birth of date error", func(t *testing.T) {
		// Setup echo
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

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "tanggal lahir"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		body := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"Format date of birth tidak sesuai (YYYY-MM-DD)",
			},
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.PutEditCustomer(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Service input validation error handling", func(t *testing.T) {
		// Setup echo
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

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		name := ""
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		body := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"Name diperlukan",
				"Name terlalu pendek",
			},
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("PutEditCustomer", body).Return(errors.Wrap(ErrInputValidation, "Name diperlukan;Name terlalu pendek"))

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.PutEditCustomer(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Service internal server error handling", func(t *testing.T) {
		// Setup echo
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

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "0812",
			Name:            "rafi ccd",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		name := "something"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		body := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("PutEditCustomer", body).Return(errors.Wrap(ErrInternalServer, "test error"))

		req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.PutEditCustomer(ctx), ctx)
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
