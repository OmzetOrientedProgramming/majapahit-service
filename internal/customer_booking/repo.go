package customerbooking

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo used to initialize repo
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetListCustomerBookingWithPagination(params ListRequest) (*List, error)
}

func (r repo) GetListCustomerBookingWithPagination(params ListRequest) (*List, error) {
	var listCustomerBooking List
	listCustomerBooking.CustomerBookings = make([]CustomerBooking, 0)
	listCustomerBooking.TotalCount = 0

	query := "SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4"
	err := r.db.Select(&listCustomerBooking.CustomerBookings, query, params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit)

	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = "SELECT COUNT(id) FROM bookings WHERE place_id = $1"
	err = r.db.Get(&listCustomerBooking.TotalCount, query, params.PlaceID)

	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &listCustomerBooking, nil
}
