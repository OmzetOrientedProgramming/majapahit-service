package place

import (
	"math"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repo interface {
	GetPlaceDetail(int) (*PlaceDetail, error)
	GetAverageRatingAndReviews(int) (*AverageRatingAndReviews, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	var result PlaceDetail

	query := "SELECT id, name, image, distance, address, description, open_hour, close_hour FROM places WHERE id = $1"
	err := r.db.Get(&result, query, placeId)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

func (r *repo) GetAverageRatingAndReviews(placeId int) (*AverageRatingAndReviews, error) {
	var result AverageRatingAndReviews
	result.Reviews = make([]UserReview, 0)

	query := "SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeId)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var sum_rating int

	query = "SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1"
	err = r.db.Get(&sum_rating, query, placeId)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var averageRating float64 = float64(sum_rating) / float64(result.ReviewCount)
	var roundedAverageRating float64 = math.Round(averageRating*100) / 100
	result.AverageRating = roundedAverageRating

	query = "SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2"
	err = r.db.Select(&result.Reviews, query, placeId)
	if err != nil {

	}

	return &result, nil
}
