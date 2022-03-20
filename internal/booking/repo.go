package booking

import (
	"github.com/jmoiron/sqlx"
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

func (r *repo) GetDetail(BookingID int) (*Detail, error) {
	var bookingDetail Detail

	query := "SELECT id, date, start_time, end_time, capacity, status, created_at FROM bookings WHERE id = $1"
	err := r.db.Get(&bookingDetail, query, BookingID)
	if err != nil {

	}

	return &bookingDetail, nil
}
