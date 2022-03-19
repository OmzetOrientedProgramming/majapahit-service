package firebaseauth

import "github.com/pkg/errors"

var (
	// ErrInternalServer for unknown error from firebase package
	ErrInternalServer = errors.New("internal server error")

	// ErrInputValidation for input validation error from firebase package
	ErrInputValidation = errors.New("input validation error")
)
