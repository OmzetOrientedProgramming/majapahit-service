package place

type PlacesList struct {
	Places     []Place `json:"places"`
	TotalCount int     `json:"total_count"`
}

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

type PlacesListRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Path  string `json:"path"`
}

type PlacesRatingAndReviewCount struct {
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"review_count" db:"review_count"'`
}
