package place

import "github.com/pkg/errors"

var (
	ErrInternalServerError  = errors.New("internal server error")
	ErrRatingNotFound       = errors.New("there is no rating yet")
	ErrInputValidationError = errors.New("input validation error")
)
