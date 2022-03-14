package auth

import "time"

// Customer represent a customer on service (business logic) layer
type Customer struct {
	ID          int
	DateOfBirth time.Time
	Gender      bool
	PhoneNumber string
	Name        string
	Status      int
}

// CustomerModel represent customer table on database, which is the child of UserModel
type CustomerModel struct {
	ID          int       `db:"id"`
	DateOfBirth time.Time `db:"date_of_birth"`
	Gender      bool      `db:"gender"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	UserID int `db:"user_id"`
}

// UserModel represent user table on database, which is the parent of CustomerModel
type UserModel struct {
	ID          int       `db:"id"`
	PhoneNumber string    `db:"phone_number"`
	Name        string    `db:"name"`
	Status      int       `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// CheckPhoneNumberRequest represent request body for checking phone number
type CheckPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// VerifyOTPRequest represent request body for verifying phone number
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

// RegisterRequest represent request body for registering
type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
}

// TwillioCredentials represent credentials for twillio service
type TwillioCredentials struct {
	SID        string
	AccountSID string
	AuthToken  string
}
