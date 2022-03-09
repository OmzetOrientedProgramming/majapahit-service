package item

import "github.com/pkg/errors"

var (
	ErrInternalServerError  = errors.New("internal server error")
	ErrInputValidationError = errors.New("input validation error")
)

