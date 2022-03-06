package place

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repo interface {
	GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error)
	GetPlaceRatingAndReviewCountByPlaceID(int) (*PlacesRatingAndReviewCount, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r repo) GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error) {
	var placeList PlacesList
	placeList.Places = make([]Place, 0)
	placeList.TotalCount = 0

	query := "SELECT id, name, description, address FROM places LIMIT $1 OFFSET $2"
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

func (r repo) GetPlaceRatingAndReviewCountByPlaceID(placeID int) (*PlacesRatingAndReviewCount, error) {
	var result PlacesRatingAndReviewCount

	query := "SELECT COUNT(rating) as review_count, COALESCE (AVG(rating), 0.0) as rating FROM reviews WHERE place_id = $1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}
