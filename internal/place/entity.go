package place

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
