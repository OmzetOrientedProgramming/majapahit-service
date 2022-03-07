package businessadminauth

import "time"

type ListUser struct {
	BusinessAdmins []BusinessAdmin `json:"business_admins"`
	Count          int             `json:"count"`
}

type BusinessAdmin struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	PlaceName   string `json:"place_name"`
	IsActive    bool   `json:"is_active"`
	BankAccount int    `json:"bank_account"`
	Balance     int    `json:"balance"`
	DateJoined  string `json:"date_joined"`
	DateUpdated string `json:"date_updated"`
}

type RegisterBusinessAdminRequest struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Capacity       int       `json:"capacity"`
	Description    string    `json:"description"`
	UserID         int       `json:"user_id"`
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
