package review

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

func (m *MockService) InsertBookingReview(review BookingReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func TestHandler_InsertBookingReview(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		placeID := 1
		bookingID := 1
		content := "Test Review"
		rating := 5

		review := BookingReview{
			UserID: userID,
			PlaceID: placeID,
			BookingID: bookingID,
			Content: content,
			Rating: rating,
		}

		payload, _ := json.Marshal(review)

		expectedResponse := util.APIResponse{
			Status:  http.StatusCreated,
			Message: "Booking review is successfully recorded.",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("InsertBookingReview", review).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/booking/:bookingID/review")
		c.SetParamNames("bookingID")
		c.SetParamValues("1")
		c.Set("userFromDatabase", &userModel)
		c.Set("userFromFirebase", &userData)

		assert.NoError(t, mockHandler.InsertBookingReview(c))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Forbidden error", func(t *testing.T) {
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

		userID := userFromDatabase.ID
		placeID := 1
		bookingID := 1
		content := "Test Review"
		rating := 5

		review := BookingReview{
			UserID: userID,
			PlaceID: placeID,
			BookingID: bookingID,
			Content: content,
			Rating: rating,
		}

		payload, _ := json.Marshal(review)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/booking/:bookingID/review")
		c.SetParamNames("bookingID")
		c.SetParamValues("1")
		c.Set("userFromDatabase", &userFromDatabase)
		c.Set("userFromFirebase", &userDataFailed)

		mockService.On("InsertBookingReview", review).Return(nil)
		util.ErrorHandler(mockHandler.InsertBookingReview(c), c)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Service internal server error", func(t *testing.T) {
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
			Name:            "rafi",
			Status:          util.StatusCustomer,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		userID := userModel.ID
		placeID := 1
		bookingID := 1
		content := "Test Review"
		rating := 5

		review := BookingReview{
			UserID: userID,
			PlaceID: placeID,
			BookingID: bookingID,
			Content: content,
			Rating: rating,
		}

		payload, _ := json.Marshal(review)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/v1/booking/:bookingID/review")
		c.SetParamNames("bookingID")
		c.SetParamValues("1")
		c.Set("userFromDatabase", &userModel)
		c.Set("userFromFirebase", &userData)

		mockService.On("InsertBookingReview", review).Return(errors.Wrap(ErrInternalServer, "test error"))

		err := mockHandler.InsertBookingReview(c)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})

	// t.Run("Input validation error", func(t *testing.T) {
	// 	// Setup echo
	// 	e := echo.New()

	// 	userData := firebaseauth.UserDataFromToken{
	// 		Kind: "",
	// 		Users: []firebaseauth.User{
	// 			{
	// 				LocalID: "",
	// 				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
	// 					{
	// 						ProviderID:  "phone",
	// 						RawID:       "",
	// 						PhoneNumber: "",
	// 						FederatedID: "",
	// 						Email:       "",
	// 					},
	// 				},
	// 				LastLoginAt:       "",
	// 				CreatedAt:         "",
	// 				PhoneNumber:       "",
	// 				LastRefreshAt:     time.Time{},
	// 				Email:             "",
	// 				EmailVerified:     false,
	// 				PasswordHash:      "",
	// 				PasswordUpdatedAt: 0,
	// 				ValidSince:        "",
	// 				Disabled:          false,
	// 			},
	// 		},
	// 	}

	// 	userModel := user.Model{
	// 		ID:              1,
	// 		PhoneNumber:     "0812",
	// 		Name:            "rafi",
	// 		Status:          util.StatusCustomer,
	// 		FirebaseLocalID: "",
	// 		Email:           "",
	// 		CreatedAt:       time.Time{},
	// 		UpdatedAt:       time.Time{},
	// 	}

	// 	userID := userModel.ID
	// 	placeID := 1
	// 	bookingID := 1
	// 	content := strings.Repeat("Test Review", 60)
	// 	rating := 6

	// 	review := BookingReview{
	// 		UserID: userID,
	// 		PlaceID: placeID,
	// 		BookingID: bookingID,
	// 		Content: content,
	// 		Rating: rating,
	// 	}

	// 	payload, _ := json.Marshal(review)

	// 	expectedResponse := util.APIResponse{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "input validation error",
	// 		Errors: []string{
	// 			"Review melebihi 500 karakter.",
	// 			"Rating invalid. Maksimum rating adalah 5.",
	// 		},
	// 	}

	// 	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// 	mockService := new(MockService)
	// 	mockHandler := NewHandler(mockService)

	// 	mockService.On("InsertBookingReview", review).Return(ErrInputValidation)

	// 	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
	// 	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.SetPath("/api/v1/booking/:bookingID/review")
	// 	c.SetParamNames("bookingID")
	// 	c.SetParamValues("1")
	// 	c.Set("userFromDatabase", &userModel)
	// 	c.Set("userFromFirebase", &userData)

	// 	assert.NoError(t, mockHandler.InsertBookingReview(c))
	// 	mockService.AssertExpectations(t)
	// 	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// 	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	// })

	// t.Run("Binding error", func(t *testing.T) {
	// 	// Setup echo
	// 	e := echo.New()

	// 	userData := firebaseauth.UserDataFromToken{
	// 		Kind: "",
	// 		Users: []firebaseauth.User{
	// 			{
	// 				LocalID: "",
	// 				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
	// 					{
	// 						ProviderID:  "phone",
	// 						RawID:       "",
	// 						PhoneNumber: "",
	// 						FederatedID: "",
	// 						Email:       "",
	// 					},
	// 				},
	// 				LastLoginAt:       "",
	// 				CreatedAt:         "",
	// 				PhoneNumber:       "",
	// 				LastRefreshAt:     time.Time{},
	// 				Email:             "",
	// 				EmailVerified:     false,
	// 				PasswordHash:      "",
	// 				PasswordUpdatedAt: 0,
	// 				ValidSince:        "",
	// 				Disabled:          false,
	// 			},
	// 		},
	// 	}

	// 	userModel := user.Model{
	// 		ID:              1,
	// 		PhoneNumber:     "0812",
	// 		Name:            "rafi",
	// 		Status:          util.StatusCustomer,
	// 		FirebaseLocalID: "",
	// 		Email:           "",
	// 		CreatedAt:       time.Time{},
	// 		UpdatedAt:       time.Time{},
	// 	}

	// 	userID := userModel.ID
	// 	placeID := 1
	// 	bookingID := 1
	// 	content := "Test Review"
	// 	exceededRating := 6

	// 	mockService := new(MockService)
	// 	mockHandler := NewHandler(mockService)



	// 	// type MockBookingReview struct {
	// 	// 	UserID		int    	`db:"user_id"`
	// 	// 	PlaceID 	int    	`db:"place_id"`
	// 	// 	BookingID 	int		`json:"booking_id" db:"booking_id"`
	// 	// 	Content		string	`json:"content" db:"content"`
	// 	// 	Rating 		int		`json:"rating" db:"rating"`
	// 	// }

	// 	review := BookingReview{
	// 		UserID: userID,
	// 		PlaceID: placeID,
	// 		BookingID: bookingID,
	// 		Content: content,
	// 		Rating: exceededRating,
	// 	}

	// 	payload, _ := json.Marshal(review)

	// 	expectedResponse := util.APIResponse{
	// 		Status:  http.StatusInternalServerError,
	// 		Message: "internal server error",
	// 	}

	// 	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// 	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
	// 	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.SetPath("/api/v1/booking/:bookingID/review")
	// 	c.SetParamNames("bookingID")
	// 	c.SetParamValues("1")
	// 	c.Set("userFromDatabase", &userModel)
	// 	c.Set("userFromFirebase", &userData)

	// 	mockService.On("InsertBookingReview", review).Return(nil)
	// 	util.ErrorHandler(mockHandler.InsertBookingReview(c), c)
	// 	mockService.AssertExpectations(t)
	// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// 	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	// })


}


