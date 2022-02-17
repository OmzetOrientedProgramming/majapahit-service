package checkup

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	response      = fmt.Sprintf("{\"status\":200,\"message\":\"application up!\"}\n")
	errorResponse = fmt.Sprintf("{\"status\":500,\"message\":\"internal server error\"}\n")
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetApplicationCheckUp() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestHandler_GetApplicationCheckUp(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/check-up")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Mock Service Expectation
	mockService.On("GetApplicationCheckUp").Return(true, nil)

	// Assert
	assert.NoError(t, h.GetApplicationCheckUp(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, response, rec.Body.String())
}

func TestHandler_GetApplicationCheckUpFailed(t *testing.T) {
	// Setup Echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/check-up")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Mock Service Expectation
	mockService.On("GetApplicationCheckUp").Return(false, errors.New("there is an error"))

	// Assert
	assert.NoError(t, h.GetApplicationCheckUp(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, errorResponse, rec.Body.String())
}
