package place

import (
	"github.com/jmoiron/sqlx"
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
	panic("Implement This Method!")
}
