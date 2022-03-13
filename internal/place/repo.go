package place

import (
	"math"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error)
	GetPlaceRatingAndReviewCountByPlaceID(int) (*PlacesRatingAndReviewCount, error)
	GetPlaceDetail(int) (*PlaceDetail, error)
	GetAverageRatingAndReviews(int) (*AverageRatingAndReviews, error)
}

type repo struct {
	db *sqlx.DB
}

// NewRepo used to initialize repo
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetPlaceDetail(placeID int) (*PlaceDetail, error) {
	var result PlaceDetail

	query := "SELECT id, name, image, distance, address, description, open_hour, close_hour, booking_price, min_slot_booking, max_slot_booking FROM places WHERE id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

func (r *repo) GetAverageRatingAndReviews(placeID int) (*AverageRatingAndReviews, error) {
	var result AverageRatingAndReviews
	result.Reviews = make([]UserReview, 0)

	query := "SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var sumRating int

	query = "SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1"
	err = r.db.Get(&sumRating, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var averageRating float64 = float64(sumRating) / float64(result.ReviewCount)
	var roundedAverageRating float64 = math.Round(averageRating*100) / 100
	result.AverageRating = roundedAverageRating

	query = "SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2"
	err = r.db.Select(&result.Reviews, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

// GetPlacesListWithPagination will do the query to database for getting list places data
func (r repo) GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error) {
	var placeList PlacesList
	placeList.Places = make([]Place, 0)
	placeList.TotalCount = 0

	query := "SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2"
	err := r.db.Select(&placeList.Places, query, params.Limit, (params.Page-1)*params.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			placeList.Places = make([]Place, 0)
			placeList.TotalCount = 0
			return &placeList, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = "SELECT COUNT(id) FROM places"
	err = r.db.Get(&placeList.TotalCount, query)
	if err != nil {
		if err == sql.ErrNoRows {
			placeList.Places = make([]Place, 0)
			placeList.TotalCount = 0
			return &placeList, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &placeList, nil
}

// GetPlaceRatingAndReviewCountByPlaceID will do the query to database for getting review and review count data
func (r repo) GetPlaceRatingAndReviewCountByPlaceID(placeID int) (*PlacesRatingAndReviewCount, error) {
	var result PlacesRatingAndReviewCount

	query := "SELECT COUNT(rating) as review_count, COALESCE (AVG(rating), 0.0) as rating FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}
