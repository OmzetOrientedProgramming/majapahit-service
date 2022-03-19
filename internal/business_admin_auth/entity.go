package businessadminauth

import "time"

// User is a media to retrieve the user_id
type User struct {
	ID          int    `json:"id" db:"id"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Status      string `json:"status" db:"status"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}

// RegisterBusinessAdminRequest is a media to bind JSON request
type RegisterBusinessAdminRequest struct {
	AdminPhoneNumber        string  `json:"admin_phone_number"`
	AdminEmail              string  `json:"admin_email"`
	AdminBankAccount        string  `json:"admin_bank_account"`
	AdminName               string  `json:"admin_name"`
	AdminBankAccountName    string  `json:"admin_bank_account_name"`
	PlaceName               string  `json:"place_name"`
	PlaceAddress            string  `json:"place_address"`
	PlaceCapacity           int     `json:"place_capacity"`
	PlaceDescription        string  `json:"place_description"`
	PlaceInterval           int     `json:"place_interval"`
	PlaceOpenHour           string  `json:"place_open_hour"`
	PlaceCloseHour          string  `json:"place_close_hour"`
	PlaceImage              string  `json:"place_image"`
	PlaceMinIntervalBooking int     `json:"place_min_interval_booking"`
	PlaceMaxIntervalBooking int     `json:"place_max_interval_booking"`
	PlaceMinSlotBooking     int     `json:"place_min_slot_booking"`
	PlaceMaxSlotBooking     int     `json:"place_max_slot_booking"`
	PlaceLat                float64 `json:"place_lat"`
	PlaceLong               float64 `json:"place_long"`
}

// LoginCredential is a media mainly to put the new-generated password
type LoginCredential struct {
	PlaceName string `json:"place_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// BusinessAdmin represent the business admin entity on business logic
type BusinessAdmin struct {
	ID                int
	Name              string
	PhoneNumber       string
	Email             string
	Password          string
	Status            int
	Balance           float64
	BankAccountNumber string
	BankAccountName   string
}

// BusinessAdminModel is a database representation of business_owners table
type BusinessAdminModel struct {
	ID                int       `db:"id"`
	Balance           float64   `db:"balance"`
	BankAccountNumber string    `db:"bank_account_number"`
	UserID            int       `db:"user_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	BankAccountName   string    `db:"bank_account_name"`
}
