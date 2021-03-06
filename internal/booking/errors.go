package booking

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used if there is error that came from the server
	ErrInternalServerError = errors.New("internal server error")

	// ErrInputValidationError is used if there is input validation error of client given data
	ErrInputValidationError = errors.New("input validation error")

	// ErrNotFound is used if resource not found
	ErrNotFound = errors.New("not found")
)
