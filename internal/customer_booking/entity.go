package customerbooking

// CustomerBooking contains information that are needed by Business Admin
type CustomerBooking struct {
	ID           	int    	`json:"id"`
	CustomerName	string 	`json:"name"`
	Capacity     	int    	`json:"capacity"`
	Date			string	`json:"date"`
	StartTime    	string 	`json:"start_time" db:"start_time"`
	EndTime      	string 	`json:"end_time" db:"end_time"`
}

// List is a container for customer bookings
type List struct {
	CustomerBooking []CustomerBooking `json:"booking_customer"`
	TotalCount      int               `json:"total_count"`
}

// ListRequest consists of request data from client
type ListRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	State int    `json:"state"`
	Path  string `json:"path"`
}
