package util

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ErrorUnwrap function will unwrap all the error and return the list of error
func ErrorUnwrap(err error) ([]string, string) {
	errString := strings.Split(err.Error(), ":")
	if len(errString) <= 1 {
		return []string{}, ""
	}
	errList, errMessage := errString[len(errString)-2], errString[len(errString)-1]

	// trim whitespace
	errMessage = strings.TrimSpace(errMessage)
	errListSplit := strings.Split(errList, ";")

	var errListSplitTrimmed []string
	for _, i := range errListSplit {
		errListSplitTrimmed = append(errListSplitTrimmed, strings.TrimSpace(i))
	}

	return errListSplitTrimmed, strings.TrimSpace(errMessage)
}

// ErrorWrapWithContext for wrapping error with context
func ErrorWrapWithContext(ctx echo.Context, statusCode int, err error, additionalMessage ...string) error {
	var message []string
	for _, i := range additionalMessage {
		message = append(message, i)
	}

	ctx.Set("errorCode", statusCode)
	if len(message) != 0 {
		return errors.Wrap(err, strings.Join(message, ";"))
	}

	return err
}

// ErrorHandler for handling HTTP error from echo
func ErrorHandler(err error, c echo.Context) {
	var (
		errorList []string
		message   string
		status    int
	)

	errorList, message = ErrorUnwrap(err)
	status = http.StatusInternalServerError

	// Handle error from echo
	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
		message = he.Message.(string)
	}

	// Handle error from internal
	statusCodeFromContext, ok := c.Get("errorCode").(int)
	if ok {
		status = statusCodeFromContext
	}

	logrus.Error(err.Error())
	switch status {
	case http.StatusInternalServerError:
		c.JSON(status, APIResponse{Status: status, Message: message})
	default:
		c.JSON(status, APIResponse{Status: status, Message: message, Errors: errorList})
	}
}
