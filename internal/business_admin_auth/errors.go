package businessadminauth

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used to mark errors on the server-side
	ErrInternalServerError = errors.New("internal server error")
	// ErrInputValidationError is used to mark error regarding of validating inputs
	ErrInputValidationError = errors.New("input validation error")
)
