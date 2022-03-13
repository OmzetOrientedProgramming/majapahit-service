package place

// Detail contain important information in Place
type Detail struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Image         string       `json:"image"`
	Distance      float64      `json:"distance"`
	Address       string       `json:"address"`
	Description   string       `json:"description"`
	OpenHour      string       `json:"open_hour" db:"open_hour"`
	CloseHour     string       `json:"close_hour" db:"close_hour"`
	BookingPrice  int          `json:"booking_price" db:"booking_price"`
	MinSlot       int          `json:"min_slot" db:"min_slot_booking"`
	MaxSlot       int          `json:"max_slot" db:"max_slot_booking"`
	AverageRating float64      `json:"average_rating" db:"rating"`
	ReviewCount   int          `json:"review_count"`
	Reviews       []UserReview `json:"reviews"`
}

// AverageRatingAndReviews contain 2 reviews, average rating, and review count of place
type AverageRatingAndReviews struct {
	AverageRating float64      `json:"average_rating"`
	ReviewCount   int          `json:"review_count" db:"count_review"`
	Reviews       []UserReview `json:"reviews"`
}

// UserReview will wrap Review of each user and another information
type UserReview struct {
	User    string  `json:"user"`
	Rating  float64 `json:"rating"`
	Content string  `json:"content"`
}

// PlacesList will wrap the Places with total count for pagination purposes
type PlacesList struct {
	Places     []Place `json:"places"`
	TotalCount int     `json:"total_count"`
}

// Place will contain the data that needed for list of place
type Place struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Distance    float64 `json:"distance"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count"`
	Image       string  `json:"image"`
}

// PlacesListRequest will wrap request data from client
type PlacesListRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Path  string `json:"path"`
}

// PlacesRatingAndReviewCount will wrap data of rating and review count
type PlacesRatingAndReviewCount struct {
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count" db:"review_count"`
}
