package customer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func (m *MockService) RetrieveCustomerProfile(userID int) (*Profile, error) {
	args := m.Called(userID)
	customerProfile := args.Get(0).(Profile)
	return &customerProfile, args.Error(1)
}

func TestHandler_RetrieveCustomerProfile(t *testing.T) {

	t.Run("success", func(t *testing.T) {
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
	
		// Setup echo
		e := echo.New()
	
		// import "net/url"
		q := make(url.Values)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)
		c.Set("userFromDatabase", &userFromDatabase)
	
		mockService := new(MockService)
		h := NewHandler(mockService)
	
		// Setting up Env
		t.Setenv("BASE_URL", "localhost:8080")
	
		customerProfile := Profile{
			PhoneNumber: 		"08123456789",
			Name:               "test_name_profile",
			Gender: 			0,
			DateOfBirth: 		time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
			ProfilePicture:     "test_image_profile",
		}
	
		expectedResponseJSON, _ := json.Marshal(util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    customerProfile,
		})
	
		mockService.On("RetrieveCustomerProfile", userFromDatabase.ID).Return(customerProfile, nil)
	
		if assert.NoError(t, h.RetrieveCustomerProfile(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
		}
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
	
		// Setup echo
		e := echo.New()
	
		// import "net/url"
		q := make(url.Values)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userDataFailed)
		c.Set("userFromDatabase", &userFromDatabase)
	
		mockService := new(MockService)
		h := NewHandler(mockService)
	
		// Setting up Env
		t.Setenv("BASE_URL", "localhost:8080")
	
		customerProfile := Profile{
			PhoneNumber: 		"08123456789",
			Name:               "test_name_profile",
			Gender: 			0,
			DateOfBirth: 		time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
			ProfilePicture:     "test_image_profile",
		}

		mockService.On("RetrieveCustomerProfile", userFromDatabase.ID).Return(customerProfile, nil)
		util.ErrorHandler(h.RetrieveCustomerProfile(c), c)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Internal server error", func(t *testing.T) {	
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

		customerProfile := Profile{
			PhoneNumber: 		"08123456789",
			Name:               "test_name_profile",
			Gender: 			0,
			DateOfBirth: 		time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
			ProfilePicture:     "test_image_profile",
		}

		q := make(url.Values)
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userFromFirebase", &userData)
		c.Set("userFromDatabase", &userFromDatabase)

		mockService := new(MockService)
		h := NewHandler(mockService)

		// Setting up Env
		t.Setenv("BASE_URL", "localhost:8080")

		mockService.On("RetrieveCustomerProfile", userFromDatabase.ID).Return(customerProfile, errors.Wrap(ErrInternalServerError, "test error"))
		util.ErrorHandler(h.RetrieveCustomerProfile(c), c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}