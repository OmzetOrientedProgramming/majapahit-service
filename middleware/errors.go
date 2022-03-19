package middleware

import "github.com/pkg/errors"

var (
	// ErrForbidden if the role is forbidden to access something
	ErrForbidden = errors.New("user does not have access to this endpoint")

	// ErrInputValidationError for the error on input validation
	ErrInputValidationError = errors.New("input validation error")
)
