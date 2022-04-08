package booking

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
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

// Repo interface for defining function that must have by repo
type Repo interface {
	GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, error)
	GetBookingData(params GetBookingDataParams) (*[]DataForCheckAvailableSchedule, error)
	GetTimeSlotsData(placeID int, selectedDate ...time.Time) (*[]TimeSlot, error)
	GetPlaceCapacity(placeID int) (*PlaceOpenHourAndCapacity, error)
	CheckedItem(ids []CheckedItemParams) (*[]CheckedItemParams, bool, error)
	CreateBookingItems(items []CreateBookingItemsParams) (*CreateBookingItemsResponse, error)
	CreateBooking(booking CreateBookingParams) (*CreateBookingResponse, error)
	UpdateTotalPrice(params UpdateTotalPriceParams) (bool, error)
	GetDetail(int) (*Detail, error)
	GetItemWrapper(int) (*ItemsWrapper, error)
	GetTicketPriceWrapper(int) (*TicketPriceWrapper, error)
	UpdateBookingStatus(int, int) error
	GetMyBookingsOngoing(localID string) (*[]Booking, error)
	GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, error)
	InsertXenditInformation(params XenditInformation) (bool, error)
	UpdateBookingStatusByXenditID(string, int) error
	GetPlaceBookingPrice(placeID int) (float64, error)
}

func (r repo) GetListCustomerBookingWithPagination(params ListRequest) (*ListBooking, error) {
	var listCustomerBooking ListBooking
	listCustomerBooking.CustomerBookings = make([]CustomerBooking, 0)
	listCustomerBooking.TotalCount = 0

	query := `SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
			FROM bookings b, users u, places p 
			WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
			ORDER BY b.date DESC LIMIT $3 OFFSET $4`
	err := r.db.Select(&listCustomerBooking.CustomerBookings, query, params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit)

	if err != nil {
		if err == sql.ErrNoRows {
			listCustomerBooking.CustomerBookings = make([]CustomerBooking, 0)
			listCustomerBooking.TotalCount = 0
			return &listCustomerBooking, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = "SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1 AND b.status = $2"
	err = r.db.Get(&listCustomerBooking.TotalCount, query, params.UserID, params.State)

	if err != nil {
		if err == sql.ErrNoRows {
			listCustomerBooking.CustomerBookings = make([]CustomerBooking, 0)
			listCustomerBooking.TotalCount = 0
			return &listCustomerBooking, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &listCustomerBooking, nil
}

func (r repo) InsertXenditInformation(params XenditInformation) (bool, error) {
	query := "UPDATE bookings SET xendit_id = $1, invoices_url = $2 WHERE id = $3"

	_, err := r.db.Exec(query, params.XenditID, params.InvoicesURL, params.BookingID)
	if err != nil {
		return false, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return true, nil
}

func (r repo) CheckedItem(ids []CheckedItemParams) (*[]CheckedItemParams, bool, error) {
	var itemsFromDatabase []CheckedItemParams
	query := "SELECT id, place_id FROM items WHERE place_id = $1"

	counter := 2
	var arguments []interface{}
	arguments = append(arguments, ids[0].PlaceID) // need to validate if all place id is the same, will do on service

	var additionalQuery []string
	for _, id := range ids {
		additionalQuery = append(additionalQuery, fmt.Sprintf("id = $%d", counter))
		arguments = append(arguments, id.ID)
		counter++
	}

	query += fmt.Sprintf(" AND (%s)", strings.Join(additionalQuery, " OR "))

	err := r.db.Select(&itemsFromDatabase, query, arguments...)
	if err != nil {
		return nil, false, errors.Wrap(ErrInternalServerError, err.Error())
	}

	if len(ids) != len(itemsFromDatabase) {
		return &itemsFromDatabase, false, errors.Wrap(ErrInputValidationError, "items not found")

	}

	return &itemsFromDatabase, true, nil
}

func (r repo) CreateBookingItems(items []CreateBookingItemsParams) (*CreateBookingItemsResponse, error) {
	query := `INSERT INTO 
					booking_items (item_id, booking_id, qty, total_price)
				VALUES`

	counter := 1
	var totalPrice float64
	var itemDataArgs []interface{}
	var itemsDataQuery []string
	for _, item := range items {
		itemsDataQuery = append(itemsDataQuery, fmt.Sprintf(" ($%d, $%d, $%d, $%d) ", counter, counter+1, counter+2, counter+3))
		itemDataArgs = append(itemDataArgs, item.ItemID)
		itemDataArgs = append(itemDataArgs, item.BookingID)
		itemDataArgs = append(itemDataArgs, item.Qty)
		itemDataArgs = append(itemDataArgs, item.TotalPrice)
		counter += 4
		totalPrice += item.TotalPrice
	}

	query += strings.Join(itemsDataQuery, ",")

	_, err := r.db.Exec(query, itemDataArgs...)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &CreateBookingItemsResponse{TotalPrice: totalPrice}, nil
}

func (r repo) CreateBooking(booking CreateBookingParams) (*CreateBookingResponse, error) {
	var bookingID CreateBookingResponse

	query := `INSERT INTO 
					bookings (user_id, place_id, date, start_time, end_time, capacity, status, total_price)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id
				`

	err := r.db.QueryRow(query, booking.UserID, booking.PlaceID, booking.Date, booking.StartTime, booking.EndTime, booking.Capacity, booking.Status, booking.TotalPrice).Scan(&bookingID.ID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &bookingID, nil
}

func (r repo) UpdateTotalPrice(params UpdateTotalPriceParams) (bool, error) {
	query := `UPDATE bookings SET total_price = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(query, params.TotalPrice, params.BookingID)
	if err != nil {
		return false, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return true, nil
}

func (r repo) GetBookingData(params GetBookingDataParams) (*[]DataForCheckAvailableSchedule, error) {
	var bookingsData []DataForCheckAvailableSchedule

	query := `SELECT id, date, start_time, end_time, capacity 
				FROM bookings 
				WHERE place_id = $1
				AND (status = $2 or status = $3)
				AND date >= $4 
				AND date <= $5`
	err := r.db.Select(&bookingsData, query, params.PlaceID, util.BookingBelumMembayar, util.BookingBerhasil, params.StartDate, params.EndDate)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &bookingsData, nil
}

func (r repo) GetTimeSlotsData(placeID int, selectedDates ...time.Time) (*[]TimeSlot, error) {
	var timeSlots []TimeSlot

	query := `SELECT id, start_time, end_time, day
				FROM time_slots 
				WHERE place_id = $1 
				`

	counter := 2
	var arguments []interface{}
	arguments = append(arguments, placeID)

	var additionalQuery []string
	for _, selectedDate := range selectedDates {
		additionalQuery = append(additionalQuery, fmt.Sprintf("day = $%d", counter))
		arguments = append(arguments, int(selectedDate.Weekday()))
		counter++
	}

	query += fmt.Sprintf("AND (%s)", strings.Join(additionalQuery, " OR "))
	query += " ORDER BY day, start_time"

	err := r.db.Select(&timeSlots, query, arguments...)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &timeSlots, nil
}

func (r repo) GetPlaceCapacity(placeID int) (*PlaceOpenHourAndCapacity, error) {
	var placeData PlaceOpenHourAndCapacity

	query := `SELECT capacity, open_hour FROM places WHERE id = $1`

	err := r.db.Get(&placeData, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &placeData, nil
}

func (r repo) GetPlaceBookingPrice(placeID int) (float64, error) {
	var bookingPrice float64

	query := `SELECT COALESCE (booking_price, 0) FROM places WHERE id  = $1`

	err := r.db.Get(&bookingPrice, query, placeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0.0, errors.Wrap(ErrNotFound, fmt.Sprintf("place with id = %d not found", placeID))
		}

		return 0.0, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return bookingPrice, nil
}

func (r *repo) GetDetail(bookingID int) (*Detail, error) {
	var bookingDetail Detail

	query := `SELECT b.id, u.name, b.date, b.start_time, b.end_time, b.capacity, b.status, b.total_price, b.created_at
			  FROM bookings b, users u
			  WHERE b.id = $1 AND b.user_id = u.id`
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

func (r *repo) UpdateBookingStatus(bookingID int, newStatus int) error {
	query := "UPDATE bookings SET status = $2 WHERE id= $1"
	_, err := r.db.Exec(query, bookingID, newStatus)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}
	return nil
}

func (r *repo) UpdateBookingStatusByXenditID(xenditID string, newStatus int) error {
	query := "UPDATE bookings SET status = $2 WHERE xendit_id= $1"
	_, err := r.db.Exec(query, xenditID, newStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Wrap(ErrNotFound, fmt.Sprintf("data with xendit_id = %s is not found", xenditID))
		}
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

func (r *repo) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	var bookingList []Booking
	bookingList = make([]Booking, 0)

	query := `
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status <= 2
		ORDER BY bookings.date asc, bookings.start_time asc
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
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
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
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
