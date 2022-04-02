package booking

// Booking contains customer booking information
type Booking struct {
	ID         int    `json:"id"`
	PlaceID    int    `json:"place_id" db:"place_id"`
	PlaceName  string `json:"place_name" db:"place_name"`
	PlaceImage string `json:"place_image" db:"place_image"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time" db:"start_time"`
	EndTime    string `json:"end_time" db:"end_time"`
	Status     int    `json:"status"`
	TotalPrice int    `json:"total_price" db:"total_price"`
}

// List contains list of customer booking information
type List struct {
	Bookings 		[]Booking	`json:"bookings"`
	TotalCount	int 			`json:"total_count"`
}

// BookingsListRequest contains request params for BookingList
type BookingsListRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Path  string `json:"path"`
}

// Detail contain required information about booking
type Detail struct {
	ID               int          `json:"id"`
	Date             string       `json:"date"`
	StartTime        string       `json:"start_time" db:"start_time"`
	EndTime          string       `json:"end_time" db:"end_time"`
	Capacity         int          `json:"capacity"`
	Status           int          `json:"status"`
	CreatedAt        string       `json:"created_at" db:"created_at"`
	TotalPrice       float64      `json:"total_price"`
	TotalPriceTicket float64      `json:"total_price_ticket"`
	TotalPriceItem   float64      `json:"total_price_item" db:"total_price"`
	Items            []ItemDetail `json:"items"`
}

// TicketPriceWrapper will consist ticket price related to place
type TicketPriceWrapper struct {
	Price float64 `db:"booking_price"`
}

// ItemsWrapper will wrap information related about item
type ItemsWrapper struct {
	Items []ItemDetail `json:"items"`
}

// ItemDetail contain required information about item
type ItemDetail struct {
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

// UpdateBookingStatusRequest represent request body for updage booking status
type UpdateBookingStatusRequest struct {
	Status int `json:"status"`
}
