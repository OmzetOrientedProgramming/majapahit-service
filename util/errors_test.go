package util

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorUnwrap(t *testing.T) {
	expectation := []string{
		"error 1", "error 2",
	}
	errMessageExpectation := errors.New("test error")

	err := errors.Wrap(errMessageExpectation, strings.Join(expectation, ";"))

	errorList, message := ErrorUnwrap(err)

	assert.Equal(t, expectation, errorList)
	assert.Equal(t, errMessageExpectation.Error(), message)
}

func TestErrorWrapWithContext(t *testing.T) {
	t.Run("error with additional message", func(t *testing.T) {
		message := []string{"test error"}
		errExpected := errors.Wrap(echo.ErrInternalServerError, strings.Join(message, ";"))
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := ErrorWrapWithContext(c, http.StatusInternalServerError, echo.ErrInternalServerError, message...)

		assert.Equal(t, errExpected.Error(), err.Error())
	})

	t.Run("error with no additional message", func(t *testing.T) {
		errExpected := errors.Wrap(echo.ErrInternalServerError, "test error")
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := ErrorWrapWithContext(c, http.StatusInternalServerError, errExpected)

		assert.Equal(t, errExpected.Error(), err.Error())
	})
}

func TestErrorHandler(t *testing.T) {
	t.Run("error with status code", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("errorCode", http.StatusBadRequest)

		errExpected := errors.Wrap(errors.New("bad request"), "test error")
		expectedResponseJSON, _ := json.Marshal(APIResponse{
			Status:  http.StatusBadRequest,
			Message: "bad request",
			Errors: []string{
				"test error",
			},
		})

		ErrorHandler(errExpected, c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("error with status code", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		errExpected := errors.Wrap(errors.New("internal server error"), "test error")
		expectedResponseJSON, _ := json.Marshal(APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		ErrorHandler(errExpected, c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))

	})

	t.Run("error from echo", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		errExpected := echo.HTTPError{
			Code:     500,
			Message:  "internal server error",
			Internal: nil,
		}

		expectedResponseJSON, _ := json.Marshal(APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})

		ErrorHandler(&errExpected, c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}
