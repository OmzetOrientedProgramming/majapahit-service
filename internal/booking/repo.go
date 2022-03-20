package booking

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetDetail(int) (*Detail, error)
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

func (r *repo) GetDetail(bookingID int) (*Detail, error) {
	var bookingDetail Detail

	query := "SELECT id, date, start_time, end_time, capacity, status, total_price, created_at FROM bookings WHERE id = $1"
	err := r.db.Get(&bookingDetail, query, bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &bookingDetail, nil
}
