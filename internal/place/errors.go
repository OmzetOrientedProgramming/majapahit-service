package place

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used if there is error that came from the server
	ErrInternalServerError = errors.New("internal server error")

	// ErrRatingNotFound is used if there is no rating yet
	ErrRatingNotFound = errors.New("there is no rating yet")

	// ErrInputValidationError is used if there is input validation error of client given data
	ErrInputValidationError = errors.New("input validation error")
)
