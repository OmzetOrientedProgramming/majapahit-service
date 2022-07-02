package middleware

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type FirebaseMockRepository struct {
	mock.Mock
}

func (f *FirebaseMockRepository) SendOTP(params firebaseauth.SendOTPParams) (*firebaseauth.SendOTPResult, error) {
	args := f.Called(params)
	return args.Get(0).(*firebaseauth.SendOTPResult), args.Error(1)
}

func (f *FirebaseMockRepository) VerifyOTP(params firebaseauth.VerifyOTPParams) (*firebaseauth.VerifyOTPResult, error) {
	args := f.Called(params)
	return args.Get(0).(*firebaseauth.VerifyOTPResult), args.Error(1)
}

func (f *FirebaseMockRepository) GetUserDataFromToken(token string) (*firebaseauth.UserDataFromToken, error) {
	args := f.Called(token)
	return args.Get(0).(*firebaseauth.UserDataFromToken), args.Error(1)
}

type UserMockRepository struct {
	mock.Mock
}

func (u *UserMockRepository) GetUserIDByLocalID(localID string) (*user.Model, error) {
	args := u.Called(localID)
	return args.Get(0).(*user.Model), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()

	firebaseRepo := new(FirebaseMockRepository)
	userRepo := new(UserMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, util.StatusCustomer)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
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

	t.Run("success", func(t *testing.T) {
		expectedRes, _ := json.Marshal(&userData)

		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)
		userRepo.On("GetUserIDByLocalID", "1").Return(&userModel, nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(expectedRes)+"\n", res.Body.String())
	})

	t.Run("failed forbidden", func(t *testing.T) {
		e.GET("/", func(c echo.Context) error {
			data, _, err := ParseUserData(c, util.StatusBusinessAdmin)
			if err != nil {
				return c.JSON(http.StatusForbidden, err)
			}
			return c.JSON(http.StatusOK, data)
		}, authMock.AuthMiddleware())

		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("invalid format", func(t *testing.T) {
		req.Header.Set(echo.HeaderAuthorization, "")
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("unsupported type", func(t *testing.T) {
		req.Header.Set(echo.HeaderAuthorization, "AUTH")
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestAuthMiddlewareUnknown(t *testing.T) {
	e := echo.New()

	firebaseRepo := new(FirebaseMockRepository)
	userRepo := new(UserMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, 3)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
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

	t.Run("success business admin", func(t *testing.T) {
		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)
		userRepo.On("GetUserIDByLocalID", "1").Return(&userModel, nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})

}

func TestAuthMiddlewareBusinessAdmin(t *testing.T) {
	e := echo.New()

	userRepo := new(UserMockRepository)
	firebaseRepo := new(FirebaseMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, util.StatusBusinessAdmin)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
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

	t.Run("success business admin", func(t *testing.T) {
		expectedRes, _ := json.Marshal(&userData)

		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)
		userRepo.On("GetUserIDByLocalID", "1").Return(&userModel, nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(expectedRes)+"\n", res.Body.String())
	})

	t.Run("failed forbidden", func(t *testing.T) {
		e.GET("/", func(c echo.Context) error {
			data, _, err := ParseUserData(c, util.StatusCustomer)
			if err != nil {
				return c.JSON(http.StatusForbidden, err)
			}
			return c.JSON(http.StatusOK, data)
		}, authMock.AuthMiddleware())

		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		assert.Equal(t, http.StatusForbidden, res.Code)
	})
}

func TestAuthMiddlewareHeaderNotProvided(t *testing.T) {
	e := echo.New()

	firebaseRepo := new(FirebaseMockRepository)
	userRepo := new(UserMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, util.StatusCustomer)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Run("header not provided", func(t *testing.T) {
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestAuthMiddlewareInternalServerError(t *testing.T) {
	e := echo.New()

	firebaseRepo := new(FirebaseMockRepository)
	userRepo := new(UserMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, util.StatusCustomer)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var userData firebaseauth.UserDataFromToken

	t.Run("internal server error from firebase repo", func(t *testing.T) {
		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, errors.Wrap(firebaseauth.ErrInternalServer, "internal server error test"))

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestAuthMiddlewareInputValidationError(t *testing.T) {
	e := echo.New()

	firebaseRepo := new(FirebaseMockRepository)
	userRepo := new(UserMockRepository)
	authMock := NewAuthMiddleware(firebaseRepo, userRepo)

	e.GET("/", func(c echo.Context) error {
		data, _, err := ParseUserData(c, util.StatusCustomer)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		return c.JSON(http.StatusOK, data)
	}, authMock.AuthMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	var userData firebaseauth.UserDataFromToken

	t.Run("input validation error from firebase repo", func(t *testing.T) {
		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, errors.Wrap(firebaseauth.ErrInputValidation, "test input validation"))

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func TestAuthMiddleware_ErrorFromGetUserData(t *testing.T) {
	t.Run("failed internal server error", func(t *testing.T) {
		e := echo.New()

		firebaseRepo := new(FirebaseMockRepository)
		userRepo := new(UserMockRepository)
		authMock := NewAuthMiddleware(firebaseRepo, userRepo)

		e.GET("/", func(c echo.Context) error {
			data, _, err := ParseUserData(c, util.StatusCustomer)
			if err != nil {
				return c.JSON(http.StatusForbidden, err)
			}
			return c.JSON(http.StatusOK, data)
		}, authMock.AuthMiddleware())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
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
		var userModel user.Model

		token := "testtoken"
		firebaseRepo.On("GetUserDataFromToken", token).Return(&userData, nil)
		userRepo.On("GetUserIDByLocalID", "1").Return(&userModel, errors.Wrap(user.ErrInternalServer, "test error"))

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()

		e.ServeHTTP(res, req)
		firebaseRepo.AssertExpectations(t)

		e.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
