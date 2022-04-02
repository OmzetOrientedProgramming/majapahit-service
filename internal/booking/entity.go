package booking

import "time"

// CustomerBooking contains information that are needed by Business Admin
type CustomerBooking struct {
	ID           int    `json:"id" db:"id"`
	CustomerName string `json:"name" db:"name"`
	Capacity     int    `json:"capacity" db:"capacity"`
	Date         string `json:"date" db:"date"`
	StartTime    string `json:"start_time" db:"start_time"`
	EndTime      string `json:"end_time" db:"end_time"`
}

// List is a container for customer bookings
type ListBooking struct {
	CustomerBookings []CustomerBooking `json:"booking_customer"`
	TotalCount       int               `json:"total_count"`
}

// ListRequest consists of request data from client
type ListRequest struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	State   int    `json:"state"`
	PlaceID int    `json:"place_id"`
	Path    string `json:"path"`
}

// DataForCheckAvailableSchedule for checking schedule in db
type DataForCheckAvailableSchedule struct {
	ID        int       `db:"id"`
	Date      time.Time `db:"date"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Capacity  int       `db:"capacity"`
}

// GetBookingDataParams parameter for check booking
type GetBookingDataParams struct {
	PlaceID   int
	StartDate time.Time
	EndDate   time.Time
	StartTime time.Time
}

// TimeSlot for time slot data
type TimeSlot struct {
	ID        int       `db:"id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Day       int       `db:"day"`
}

// TimeSlotAPIResponse for wrapping response for time slots
type TimeSlotAPIResponse struct {
	ID        int    `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Day       int    `json:"day"`
}

// AvailableTimeResponse for response after check available time
type AvailableTimeResponse struct {
	Time  string `json:"time"`
	Total int    `json:"total"`
}

// AvailableDateResponse for response after check available date
type AvailableDateResponse struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

// GetAvailableTimeParams parameter for get available time
type GetAvailableTimeParams struct {
	PlaceID      int
	SelectedDate time.Time
	StartTime    time.Time
	BookedSlot   int
}

// GetAvailableDateParams parameter for get available date
type GetAvailableDateParams struct {
	PlaceID    int
	StartDate  time.Time
	Interval   int
	BookedSlot int
}

// PlaceOpenHourAndCapacity for place open hour and capacity data
type PlaceOpenHourAndCapacity struct {
	OpenHour time.Time `db:"open_hour"`
	Capacity int       `db:"capacity"`
}

// CheckedItemParams is parameter for checked item
type CheckedItemParams struct {
	ID      int `json:"id" db:"id"`
	PlaceID int `json:"place_id" db:"place_id"`
}

// CheckedItemResponse is response for checked item
type CheckedItemResponse struct {
	ID      int     `json:"id" db:"id"`
	Price   float64 `json:"price" db:"price"`
	PlaceID int     `json:"place_id" db:"place_id"`
}

// CreateBookingItemsParams for inserting to booking item table
type CreateBookingItemsParams struct {
	BookingID  int     `json:"booking_id"`
	ItemID     int     `json:"id"`
	TotalPrice float64 `json:"total_price"`
	Qty        int     `json:"qty"`
}

// CreateBookingItemsResponse struct for the response after inserting booking item data to db
type CreateBookingItemsResponse struct {
	TotalPrice float64 `json:"total_price" db:"total_price"`
}

// CreateBookingParams for inserting to booking table
type CreateBookingParams struct {
	UserID     int       `json:"user_id" db:"user_id"`
	PlaceID    int       `json:"place_id" db:"place_id"`
	Date       time.Time `json:"date"`
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
	Capacity   int       `json:"capacity"`
	Status     int       `json:"status"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
}

// CreateBookingResponse struct for the response after inserting booking data to db
type CreateBookingResponse struct {
	ID int `json:"id" db:"id"`
}

// UpdateTotalPriceParams struct for update total price of booking function params
type UpdateTotalPriceParams struct {
	BookingID  int     `json:"booking_id"`
	TotalPrice float64 `json:"total_price"`
}

// CreateBookingServiceRequest request for create booking
type CreateBookingServiceRequest struct {
	Items               []Item    `json:"items"`
	Date                time.Time `json:"date"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Count               int       `json:"count"`
	PlaceID             int       `json:"place_id"`
	UserID              int       `json:"user_id"`
	CustomerName        string    `json:"customer_name"`
	CustomerPhoneNumber string    `json:"customer_phone_number"`
}

// Item object on create booking request
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

// CreateBookingServiceResponse response for create booking
type CreateBookingServiceResponse struct {
	XenditID   string `json:"xendit_id"`
	BookingID  int    `json:"booking_id"`
	PaymentURL string `json:"payment_url"`
}

// CreateBookingRequestBody for API request body
type CreateBookingRequestBody struct {
	Items     []Item `json:"items"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Count     int    `json:"count"`
}

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
	Bookings   []Booking `json:"bookings"`
	TotalCount int       `json:"total_count"`
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
