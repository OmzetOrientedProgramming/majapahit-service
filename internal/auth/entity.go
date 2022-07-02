package auth

import "time"

// Customer represent a customer on service (business logic) layer
type Customer struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Status      string    `json:"status"`
	LocalID     string    `json:"local_id"`
}

// CustomerModel represent customer table on database, which is the child of UserModel
type CustomerModel struct {
	ID          int       `db:"id"`
	DateOfBirth time.Time `db:"date_of_birth"`
	Gender      int       `db:"gender"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	UserID      int       `db:"user_id"`
}

// UserModel represent user table on database, which is the parent of CustomerModel
type UserModel struct {
	ID              int       `db:"id"`
	PhoneNumber     string    `db:"phone_number"`
	Name            string    `db:"name"`
	Status          int       `db:"status"`
	FirebaseLocalID string    `db:"firebase_local_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Email           string    `db:"email"`
	Password        string    `db:"password"`
	Image           string    `db:"image"`
}

// CheckPhoneNumberRequest represent request body for checking phone number
type CheckPhoneNumberRequest struct {
	PhoneNumber    string `json:"phone_number"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// VerifyOTPRequest represent request body for verifying phone number
type VerifyOTPRequest struct {
	SessionInfo string `json:"session_info"`
	OTP         string `json:"otp"`
}

// RegisterRequest represent request body for registering
type RegisterRequest struct {
	FullName string `json:"full_name"`
}

// CheckPhoneNumberResponse response of check phone number endpoint
type CheckPhoneNumberResponse struct {
	SessionInfo string `json:"session_info"`
}

// VerifyOTPResult Result of verify OTP endpoint
type VerifyOTPResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	LocalID      string `json:"local_id"`
	IsNewUser    bool   `json:"is_new_user"`
	PhoneNumber  string `json:"phone_number"`
}
