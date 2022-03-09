package businessadminauth

import (
	"time"
)

type ListUser struct {
	BusinessAdmins []BusinessAdmin `json:"business_admins"`
	Count          int             `json:"count"`
}

type ListPlace struct {
	Places []Place `json:"places"`
	Count  int     `json:"count"`
}

type BusinessAdmin struct {
	ID          int       `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	PlaceID     int       `json:"place_id"`
	IsActive    bool      `json:"is_active"`
	BankAccount int       `json:"bank_account"`
	Balance     int       `json:"balance"`
	JoinedAt    time.Time `json:"date_joined"`
	UpdatedAt   time.Time `json:"date_updated"`
	Password    string    `json:"password"`
}

type Place struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Capacity       int       `json:"capacity"`
	Description    string    `json:"description"`
	AdminID        int       `json:"admin_id"`
	Interval       string    `json:"interval"`
	OpenHour       string    `json:"open_hour"`
	CloseHour      string    `json:"close_hour"`
	Image          string    `json:"image"`
	MinHourBooking int       `json:"min_hour_booking"`
	MaxHourBooking int       `json:"max_hour_booking"`
	MinSlotBooking int       `json:"min_slot_booking"`
	MaxSlotBooking int       `json:"max_slot_booking"`
	Lat            float64   `json:"lat"`
	Long           float64   `json:"long"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisterBusinessAdminRequest struct {
	AdminPhoneNumber    string  `json:"admin_phone_number"`
	AdminEmail          string  `json:"admin_email"`
	AdminBankAccount    string  `json:"admin_bank_account"`
	PlaceName           string  `json:"place_name"`
	PlaceAddress        string  `json:"place_address"`
	PlaceCapacity       int     `json:"place_capacity"`
	PlaceDescription    string  `json:"place_description"`
	PlaceInterval       string  `json:"place_interval"`
	PlaceOpenHour       string  `json:"place_open_hour"`
	PlaceCloseHour      string  `json:"place_close_hour"`
	PlaceImage          string  `json:"place_image"`
	PlaceMinHourBooking int     `json:"place_min_hour_booking"`
	PlaceMaxHourBooking int     `json:"place_max_hour_booking"`
	PlaceMinSlotBooking int     `json:"place_min_slot_booking"`
	PlaceMaxSlotBooking int     `json:"place_max_slot_booking"`
	PlaceLat            float64 `json:"place_lat"`
	PlaceLong           float64 `json:"place_long"`
}
