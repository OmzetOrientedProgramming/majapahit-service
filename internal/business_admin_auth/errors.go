package businessadminauth

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used to mark errors on the server-side
	ErrInternalServerError = errors.New("internal server error")
	// ErrInputValidationError is used to mark error regarding of validating inputs
	ErrInputValidationError = errors.New("input validation error")
	// ErrUnauthorized is used to mark error regarding of unauthorized access
	ErrUnauthorized = errors.New("unauthorized")
	// ErrNotFound is used to mark error regarding of not found resource(s)
	ErrNotFound = errors.New("not found")
)
