package customer

import "github.com/pkg/errors"

var (
	// ErrInputValidation is returned when the input is invalid
	ErrInputValidation = errors.New("input validation error")
	// ErrInternalServer is returned when the server encounters an internal error
	ErrInternalServer = errors.New("internal server error")
)
