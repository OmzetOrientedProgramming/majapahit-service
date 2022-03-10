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
