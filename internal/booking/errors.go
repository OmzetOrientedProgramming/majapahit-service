package booking

import "github.com/pkg/errors"

var (
	// ErrInternalServerError is used if there is error that came from the server
	ErrInternalServerError = errors.New("internal server error")
)
