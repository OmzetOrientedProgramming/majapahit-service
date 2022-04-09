package businessadmin

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
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetBalanceDetail(userID int) (*BalanceDetail, error) {
	args := m.Called(userID)
	ret := args.Get(0).(*BalanceDetail)
	return ret, args.Error(1)
}

func TestHandler_GetBalanceDetailSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

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

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	balanceDetail := BalanceDetail{
		LatestDisbursementDate: "27 Januari 2022",
		Balance:                2500000,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    balanceDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, nil)

	// Tes
	if assert.NoError(t, h.GetBalanceDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetBalanceDetailParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

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

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Excpectation
	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, nil)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetBalanceDetailInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

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

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetBalanceDetailBadRequestFromService(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

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

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInputValidationError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"test",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}
