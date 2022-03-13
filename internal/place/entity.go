package place

type PlaceDetail struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Image         string       `json:"image"`
	Distance      float64      `json:"distance"`
	Address       string       `json:"address"`
	Description   string       `json:"description"`
	OpenHour      string       `json:"open_hour" db:"open_hour"`
	CloseHour     string       `json:"close_hour" db:"close_hour"`
	AverageRating float64      `json:"average_rating" db:"rating"`
	ReviewCount   int          `json:"review_count"`
	Reviews       []UserReview `json:"reviews"`
}

type AverageRatingAndReviews struct {
	AverageRating float64      `json:"average_rating"`
	ReviewCount   int          `json:"review_count" db:"count_review"`
	Reviews       []UserReview `json:"reviews"`
}

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
