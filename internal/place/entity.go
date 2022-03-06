package place

import (
	"time"
)

type PlaceDetail struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Image         string       `json:"image"`
	Distance      float64      `json:"distance"`
	Address       string       `json:"address"`
	Description   string       `json:"description"`
	OpenHour      time.Time    `json:"open_hour"`
	CloseHour     time.Time    `json:"close_hour"`
	AverageRating float64      `json:"average_rating"`
	ReviewCount   int          `json:"review_count"`
	Reviews       []UserReview `json:"reviews"`
}

type UserReview struct {
	User    string  `json:"user"`
	Rating  float64 `json:"rating"`
	Content string  `json:"content"`
}
