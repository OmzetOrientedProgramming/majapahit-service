package firebaseauth

import "time"

// SendOTPParams is parameter for sending OTP
type SendOTPParams struct {
	PhoneNumber    string `json:"phone_number"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// SendOTPResult is result after calling SendOTP function
type SendOTPResult struct {
	SessionInfo string `json:"session_info"`
}

// ErrorInstanceMessage is message object of error instance from firebase
type ErrorInstanceMessage struct {
	Message string `json:"message"`
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
}

// ErrorInstance is error object from firebase
type ErrorInstance struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Errors  []ErrorInstanceMessage `json:"errors"`
}

// ErrorFromFirebase is error object from firebase
type ErrorFromFirebase struct {
	Error ErrorInstance `json:"error"`
}

// VerifyOTPParams is parameter for calling verify OTP function
type VerifyOTPParams struct {
	SessionInfo string `json:"session_info"`
	Code        string `json:"code"`
}

// VerifyOTPResult is result of calling verify OTP function
type VerifyOTPResult struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	IsNewUser    bool   `json:"isNewUser"`
	PhoneNumber  string `json:"phoneNumber"`
}

// UserDataFromToken is user information after calling get user data from token function
type UserDataFromToken struct {
	Kind  string `json:"kind"`
	Users []User `json:"users"`
}

// User object from firebase
type User struct {
	LocalID           string             `json:"localId"`
	ProviderUserInfo  []ProviderUserInfo `json:"providerUserInfo"`
	LastLoginAt       string             `json:"lastLoginAt"`
	CreatedAt         string             `json:"createdAt"`
	PhoneNumber       string             `json:"phoneNumber"`
	LastRefreshAt     time.Time          `json:"lastRefreshAt"`
	Email             string             `json:"email"`
	EmailVerified     bool               `json:"emailVerified"`
	PasswordHash      string             `json:"passwordHash"`
	PasswordUpdatedAt int                `json:"passwordUpdatedAt"`
	ValidSince        string             `json:"validSince"`
	Disabled          bool               `json:"disabled"`
}

// ProviderUserInfo object from firebase
type ProviderUserInfo struct {
	ProviderID  string `json:"providerId"`
	RawID       string `json:"rawId"`
	PhoneNumber string `json:"phoneNumber"`
	FederatedID string `json:"federatedId"`
	Email       string `json:"email"`
}
