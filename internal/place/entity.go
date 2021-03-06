package place

import "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"

// Detail contain important information in Place
type Detail struct {
	ID                 int          `json:"id"`
	Name               string       `json:"name"`
	Image              string       `json:"image"`
	Address            string       `json:"address"`
	Description        string       `json:"description"`
	OpenHour           string       `json:"open_hour" db:"open_hour"`
	CloseHour          string       `json:"close_hour" db:"close_hour"`
	BookingPrice       int          `json:"booking_price" db:"booking_price"`
	MinSlot            int          `json:"min_slot" db:"min_slot_booking"`
	MaxSlot            int          `json:"max_slot" db:"max_slot_booking"`
	Capacity           int          `json:"capacity" db:"capacity"`
	MinIntervalBooking int          `json:"min_interval_booking" db:"min_interval_booking"`
	MaxIntervalBooking int          `json:"max_interval_booking" db:"max_interval_booking"`
	AverageRating      float64      `json:"average_rating" db:"rating"`
	ReviewCount        int          `json:"review_count"`
	Reviews            []UserReview `json:"reviews"`
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
	TotalCount int     `json:"total_count" db:"total_count"`
}

// Place will contain the data that needed for list of place
type Place struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	Distance    int     `json:"distance"`
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count" db:"review_count"`
	Image       string  `json:"image"`
}

// PlacesListRequest will wrap request data from client
type PlacesListRequest struct {
	Limit int    `query:"limit"`
	Page  int    `query:"page"`
	Path  string `json:"path"`

	// Distance param
	Latitude  float64 `query:"lat"`
	Longitude float64 `query:"lng"`

	// Filter params
	Price  []string `query:"price"`
	People []string `query:"people"`
	Rating []int    `query:"rating"`

	// Sort params
	Sort string `query:"sort"`

	// Category params
	Category string `query:"category"`
}

// PlacesListResponse will wrap response data to client
type PlacesListResponse struct {
	Places     *PlacesList      `json:"places"`
	Pagination *util.Pagination `json:"pagination"`
}

// PlacesRatingAndReviewCount will wrap data of rating and review count
type PlacesRatingAndReviewCount struct {
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count" db:"review_count"`
}

// Review consist of informations about review and rating from customer
type Review struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
	Date    string `json:"created_at" db:"created_at"`
}

// ListReview is a container for review
type ListReview struct {
	Reviews    []Review `json:"review"`
	TotalCount int      `json:"total_count"`
}

// ListReviewRequest consist of request for pagination and sorting purpose
type ListReviewRequest struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Path    string `json:"path"`
	PlaceID int    `json:"place_id"`
	Rating  bool   `json:"rating"`
	Latest  bool   `json:"latest"`
}
