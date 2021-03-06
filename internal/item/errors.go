package item

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used if there is error that came from the server
	ErrInternalServerError = errors.New("internal server error")

	// ErrInputValidationError is used if there is input validation error of client given data
	ErrInputValidationError = errors.New("input validation error")

	// ErrNotFound is used if there is no item can be found
	ErrNotFound = errors.New("item not found")
)
