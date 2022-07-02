package upload

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

func (m *MockService) UploadProfilePicture(params FileRequest) (*FileResponse, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*FileResponse), args.Error(1)
}

func TestHandler_UploadProfilePicture(t *testing.T) {
	t.Run("Success Upload Profile Picture", func(t *testing.T) {
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

		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		name := userModel.Name

		body := FileRequest{
			File:         file,
			CustomerName: name,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusCreated,
			Message: "Successfully Uploaded File!",
			Data: FileResponse{
				URL: "https://res.cloudinary.com/wave-ppl/image/upload/v1650884844/Profile%20Picture/Mario%20Ganteng-Profile-Picture.png",
			},
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UploadProfilePicture", body).Return(&FileResponse{
			URL: "https://res.cloudinary.com/wave-ppl/image/upload/v1650884844/Profile%20Picture/Mario%20Ganteng-Profile-Picture.png",
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/upload/profile-picture", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		assert.NoError(t, mockHandler.UploadProfilePicture(ctx))
		mockService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rec.Code)
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

		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		name := userModelFailed.Name

		body := FileRequest{
			File:         file,
			CustomerName: name,
		}

		payload, _ := json.Marshal(body)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UploadProfilePicture", body).Return(&FileResponse{
			URL: "https://res.cloudinary.com/wave-ppl/image/upload/v1650884844/Profile%20Picture/Mario%20Ganteng-Profile-Picture.png",
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/upload/profile-picture", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModelFailed)
		ctx.Set("userFromFirebase", &userDataFailed)

		util.ErrorHandler(mockHandler.UploadProfilePicture(ctx), ctx)
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
		type MismatchFileRequest struct {
			File         int `json:"file"`
			CustomerName string
		}

		file := 1
		name := userModel.Name

		body := MismatchFileRequest{
			File:         file,
			CustomerName: name,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/upload/profile-picture", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.UploadProfilePicture(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Service error input validation", func(t *testing.T) {
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

		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		name := userModel.Name

		body := FileRequest{
			File:         file,
			CustomerName: name,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors: []string{
				"service test",
			},
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UploadProfilePicture", body).Return(nil, errors.Wrap(ErrInputValidation, "service test"))

		req := httptest.NewRequest(http.MethodPost, "/upload/profile-picture", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.UploadProfilePicture(ctx), ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("Service error status internal server error", func(t *testing.T) {
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

		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		name := userModel.Name

		body := FileRequest{
			File:         file,
			CustomerName: name,
		}

		payload, _ := json.Marshal(body)

		expectedResponse := util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}

		expectedResponseJSON, _ := json.Marshal(expectedResponse)

		mockService := new(MockService)
		mockHandler := NewHandler(mockService)

		mockService.On("UploadProfilePicture", body).Return(nil, errors.Wrap(ErrInternalServer, "service test"))

		req := httptest.NewRequest(http.MethodPost, "/upload/profile-picture", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		util.ErrorHandler(mockHandler.UploadProfilePicture(ctx), ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
