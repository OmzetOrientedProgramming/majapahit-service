package businessadminauth

import "github.com/pkg/errors"

var (
	ErrInternalServerError     = errors.New("internal server error")
	ErrInputValidationError    = errors.New("input validation error")
	ErrUserNotFound            = errors.New("user not found")
	ErrPhoneNumberAlreadyExist = errors.New("phone_number is already taken")
	ErrNameIsTooShort          = errors.New("name is too short")
)
