package cloudinary

import "github.com/pkg/errors"

var (
	// ErrInternalServer for unknown error from cloudinary package
	ErrInternalServer = errors.New("internal server error")
)
