package user

import "github.com/pkg/errors"

var (
	// ErrNotFound is returned when user is not found
	ErrNotFound = errors.New("user not found")
	// ErrInternalServer is returned when the server encounters an internal error
	ErrInternalServer = errors.New("internal server error")
)
