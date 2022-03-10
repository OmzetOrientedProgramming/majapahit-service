package place

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repo interface {
	GetPlaceDetail(placeId int) (*PlaceDetail, error)
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
