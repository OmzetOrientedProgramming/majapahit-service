package businessadminauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Login(email, password, recaptchaToken string) (*BusinessAdmin, string, error) {
	//TODO implement me
	panic("implement me")
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
