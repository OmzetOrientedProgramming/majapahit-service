package review

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo PostgreSQL for checkup module
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
	InsertBookingReview(review BookingReview) error
	RetrievePlaceID(bookingID int) (*int, error)
	CheckBookingStatus(bookingID int) (bool, error)
	UpdateBookingStatus(bookingID int) error
}

func (r repo) InsertBookingReview(review BookingReview) error {
	query := `
    INSERT INTO reviews (user_id, place_id, booking_id, content, rating)
    VALUES (:user_id, :place_id, :booking_id, :content, :rating);
	`

	_, err := r.db.NamedExec(query, review)
	if err != nil {
		return errors.Wrap(ErrInternalServer, err.Error())
	}

	return nil
}

func (r repo) RetrievePlaceID(bookingID int) (*int, error) {
	query := `
    SELECT place_id FROM bookings WHERE id=$1 LIMIT 1;
	`

	var placeID int
	err := r.db.Get(&placeID, query, bookingID)

	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}
	return &placeID, nil
}

func (r repo) CheckBookingStatus(bookingID int) (bool, error) {
	query := `
    SELECT status FROM bookings WHERE id=$1 LIMIT 1;
	`
	var bookingStatus int
	err := r.db.Get(&bookingStatus, query, bookingID)

	if err != nil {		
		return false, errors.Wrap(ErrInputValidation, "Booking tidak ditemukan")
	}

	if bookingStatus != 3 {
		return false, nil
	}

	return true, nil
}

func (r repo) UpdateBookingStatus(bookingID int) error {
	query := `
    UPDATE bookings
	SET status = 5
	WHERE id = $1
	`

	_, err := r.db.Exec(query, bookingID)
	if err != nil {
		return errors.Wrap(ErrInternalServer, err.Error())
	}

	return nil
}
