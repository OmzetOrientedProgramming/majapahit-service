package booking

import "github.com/jmoiron/sqlx"

type Repo interface {
	GetDetail(int) (*Detail, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetDetail(BookingID int) (*Detail, error) {
	panic("Implement This!")
}
