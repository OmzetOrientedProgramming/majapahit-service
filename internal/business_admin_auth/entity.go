package businessadminauth

type User struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

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

type RegisterBusinessAdminResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    LoginCredential `json:"data,omitempty"`
	Errors  []string        `json:"errors,omitempty"`
}

type LoginCredential struct {
	PlaceName string `json:"place_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
