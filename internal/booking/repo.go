package booking

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetDetail(int) (*Detail, error)
	GetItemWrapper(int) (*ItemsWrapper, error)
	GetTicketPriceWrapper(int) (*TicketPriceWrapper, error)

	GetMyBookingsOngoing(localID string) (*[]Booking, error)
	GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, error)
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

func (r *repo) GetItemWrapper(bookingID int) (*ItemsWrapper, error) {
	var bookingItems ItemsWrapper
	bookingItems.Items = make([]ItemDetail, 0)

	query := "SELECT items.name as name, items.image as image, booking_items.qty as qty, items.price as price FROM items INNER JOIN booking_items ON items.id = booking_items.item_id WHERE booking_items.booking_id = $1"
	err := r.db.Select(&bookingItems.Items, query, bookingID)

	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &bookingItems, nil
}

func (r *repo) GetTicketPriceWrapper(bookingID int) (*TicketPriceWrapper, error) {
	var ticketPrice TicketPriceWrapper

	query := "SELECT booking_price FROM places INNER JOIN bookings ON bookings.place_id = places.id WHERE bookings.id= $1"
	err := r.db.Get(&ticketPrice, query, bookingID)

	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &ticketPrice, nil
}


func (r *repo) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	var bookingList []Booking
	bookingList = make([]Booking, 0)

	query := `
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time > now()
		ORDER BY bookings.end_time DESC
	`
	err := r.db.Select(&bookingList, query, localID)
	if err != nil {
		if err == sql.ErrNoRows {
			bookingList = make([]Booking, 0)
			return &bookingList, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}
	
	return &bookingList, nil
}

func (r repo) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, error) {
	var myBookingsPrevious List
	myBookingsPrevious.Bookings = make([]Booking, 0)

	query := `
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`
	err := r.db.Select(&myBookingsPrevious.Bookings, query, localID, params.Limit, (params.Page-1)*params.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			myBookingsPrevious.Bookings = make([]Booking, 0)
			myBookingsPrevious.TotalCount = 0
			return &myBookingsPrevious, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = `
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now()
	`
	err = r.db.Get(&myBookingsPrevious.TotalCount, query, localID)
	if err != nil {
		if err == sql.ErrNoRows {
			myBookingsPrevious.Bookings = make([]Booking, 0)
			myBookingsPrevious.TotalCount = 0
			return &myBookingsPrevious, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &myBookingsPrevious, nil
}
