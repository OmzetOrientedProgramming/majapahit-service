package review

// BookingReview is a struct to accomodate inserting booking reviews
type BookingReview struct {	
	UserID		int    	`db:"user_id"`
	PlaceID 	int    	`db:"place_id"`
	BookingID 	int		`json:"booking_id" db:"booking_id"`
	Content		string	`json:"content" db:"content"`
	Rating 		int		`json:"rating" db:"rating"`
}