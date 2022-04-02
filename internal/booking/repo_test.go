package booking

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetDetailSuccess(t *testing.T) {
	bookingID := 1
	createdAtRow := time.Date(2021, time.Month(10), 26, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	bookingDetailExpected := &Detail{
		ID:             1,
		Date:           "27 Oktober 2021",
		StartTime:      "19:00",
		EndTime:        "20:00",
		Capacity:       10,
		Status:         1,
		TotalPriceItem: 100000.0,
		CreatedAt:      createdAtRow,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "date", "start_time", "end_time", "capacity", "status", "total_price", "created_at"}).
		AddRow(
			bookingDetailExpected.ID,
			bookingDetailExpected.Date,
			bookingDetailExpected.StartTime,
			bookingDetailExpected.EndTime,
			bookingDetailExpected.Capacity,
			bookingDetailExpected.Status,
			bookingDetailExpected.TotalPriceItem,
			bookingDetailExpected.CreatedAt,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date, start_time, end_time, capacity, status, total_price, created_at FROM bookings WHERE id = $1")).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	bookingDetailRetrieved, err := repoMock.GetDetail(bookingID)
	assert.Equal(t, bookingDetailExpected, bookingDetailRetrieved)
	assert.NotNil(t, bookingDetailRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetDetailInternalServerError(t *testing.T) {
	bookingID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, date, start_time, end_time, capacity, status, total_price, created_at FROM bookings WHERE id = $1")).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	bookingDetailRetrieved, err := repoMock.GetDetail(bookingID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, bookingDetailRetrieved)
}

func TestRepo_GetItemWrapperSucces(t *testing.T) {
	bookingID := 1
	itemWrapperExpected := &ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "Jus Mangga Asyik",
				Image: "ini_link_gambar_1",
				Qty:   10,
				Price: 10000.0,
			},
			{
				Name:  "Pizza with Pinapple Large",
				Image: "ini_link_gambar_2",
				Qty:   2,
				Price: 150000.0,
			},
		},
	}

	// Initialized Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Setup Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"name", "image", "qty", "price"}).
		AddRow(
			itemWrapperExpected.Items[0].Name,
			itemWrapperExpected.Items[0].Image,
			itemWrapperExpected.Items[0].Qty,
			itemWrapperExpected.Items[0].Price,
		).
		AddRow(
			itemWrapperExpected.Items[1].Name,
			itemWrapperExpected.Items[1].Image,
			itemWrapperExpected.Items[1].Qty,
			itemWrapperExpected.Items[1].Price,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT items.name as name, items.image as image, booking_items.qty as qty, items.price as price FROM items INNER JOIN booking_items ON items.id = booking_items.item_id WHERE booking_items.booking_id = $1")).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	itemWrapperRetrieved, err := repoMock.GetItemWrapper(bookingID)
	assert.Equal(t, itemWrapperExpected, itemWrapperRetrieved)
	assert.NotNil(t, itemWrapperRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetItemWrapperInternalServerError(t *testing.T) {
	bookingID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT items.name as name, items.image as image, booking_items.qty as qty, items.price as price FROM items INNER JOIN booking_items ON items.id = booking_items.item_id WHERE booking_items.booking_id = $1")).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	itemWrapperRetrieved, err := repoMock.GetItemWrapper(bookingID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, itemWrapperRetrieved)
}

func TestRepo_GetTicketPriceWrapperSuccess(t *testing.T) {
	bookingID := 1
	ticketPriceWrapperExpected := &TicketPriceWrapper{
		Price: 10000,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"booking_price"}).
		AddRow(
			ticketPriceWrapperExpected.Price,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT booking_price FROM places INNER JOIN bookings ON bookings.place_id = places.id WHERE bookings.id= $1")).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	ticketPriceWrapperRetrieved, err := repoMock.GetTicketPriceWrapper(bookingID)
	assert.Equal(t, ticketPriceWrapperExpected, ticketPriceWrapperRetrieved)
	assert.NotNil(t, ticketPriceWrapperRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetTicketPriceWrapperInternalServerError(t *testing.T) {
	bookingID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT booking_price FROM places INNER JOIN bookings ON bookings.place_id = places.id WHERE bookings.id= $1")).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	ticketPriceWrapperRetrieved, err := repoMock.GetTicketPriceWrapper(bookingID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, ticketPriceWrapperRetrieved)
}


func TestRepo_GetMyBookingsOngoingSuccess(t *testing.T) {
	localID := "abc"
	myBookingsOngoingExpected := []Booking{
		{
			ID:         1,
			PlaceID:    2,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       "2022-04-10",
			StartTime:  "08:00",
			EndTime:    "10:00",
			Status:     0,
			TotalPrice: 10000,
		}, 
		{
			ID:         2,
			PlaceID:    3,
			PlaceName:  "test_place_name",
			PlaceImage: "test_place_image",
			Date:       "2022-04-11",
			StartTime:  "09:00",
			EndTime:    "11:00",
			Status:     0,
			TotalPrice: 20000,
		},
	}
	

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "place_id", "place_name", "place_image", "date", "start_time", "end_time", "status", "total_price"}).
		AddRow(
			myBookingsOngoingExpected[0].ID,
			myBookingsOngoingExpected[0].PlaceID,
			myBookingsOngoingExpected[0].PlaceName,
			myBookingsOngoingExpected[0].PlaceImage,
			myBookingsOngoingExpected[0].Date,
			myBookingsOngoingExpected[0].StartTime,
			myBookingsOngoingExpected[0].EndTime,
			myBookingsOngoingExpected[0].Status,
			myBookingsOngoingExpected[0].TotalPrice,
		).
		AddRow(
			myBookingsOngoingExpected[1].ID,
			myBookingsOngoingExpected[1].PlaceID,
			myBookingsOngoingExpected[1].PlaceName,
			myBookingsOngoingExpected[1].PlaceImage,
			myBookingsOngoingExpected[1].Date,
			myBookingsOngoingExpected[1].StartTime,
			myBookingsOngoingExpected[1].EndTime,
			myBookingsOngoingExpected[1].Status,
			myBookingsOngoingExpected[1].TotalPrice,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time > now()
	`)).
		WithArgs(localID).
		WillReturnRows(rows)

	// Test
	myBookingsOngoingRetrieve, err := repoMock.GetMyBookingsOngoing(localID)
	assert.Equal(t, &myBookingsOngoingExpected, myBookingsOngoingRetrieve)
	assert.NotNil(t, myBookingsOngoingRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetMyBookingsOngoingEmpty(t *testing.T) {
	localID := "abc"
	myBookingsOngoingExpected := make([]Booking, 0)

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "place_id", "place_name", "place_image", "date", "start_time", "end_time", "status", "total_price"})

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time > now()
	`)).
		WithArgs(localID).
		WillReturnRows(rows)

	// Test
	myBookingsOngoingRetrieve, err := repoMock.GetMyBookingsOngoing(localID)
	assert.Equal(t, &myBookingsOngoingExpected, myBookingsOngoingRetrieve)
	assert.NotNil(t, myBookingsOngoingRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetMyBookingsOngoingInternalServerError(t *testing.T) {
	localID := "abc"

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time > now()
	`)).
		WithArgs(localID).
		WillReturnError(sql.ErrTxDone)

	// Test
	placeDetailRetrieve, err := repoMock.GetMyBookingsOngoing(localID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailRetrieve)
}

func TestRepo_GetMyBookingsPreviousWithPaginationSuccess(t *testing.T) {
	myBookingsPreviousExpected := &List{
		Bookings: []Booking{
			{
				ID:         1,
				PlaceID:    2,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-10",
				StartTime:  "08:00",
				EndTime:    "10:00",
				Status:     0,
				TotalPrice: 10000,
			}, 
			{
				ID:         2,
				PlaceID:    3,
				PlaceName:  "test_place_name",
				PlaceImage: "test_place_image",
				Date:       "2022-04-11",
				StartTime:  "09:00",
				EndTime:    "11:00",
				Status:     0,
				TotalPrice: 20000,
			},
		},
		TotalCount: 10,
	}
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
	NewRows([]string{"id", "place_id", "place_name", "place_image", "date", "start_time", "end_time", "status", "total_price"}).
	AddRow(
		myBookingsPreviousExpected.Bookings[0].ID,
		myBookingsPreviousExpected.Bookings[0].PlaceID,
		myBookingsPreviousExpected.Bookings[0].PlaceName,
		myBookingsPreviousExpected.Bookings[0].PlaceImage,
		myBookingsPreviousExpected.Bookings[0].Date,
		myBookingsPreviousExpected.Bookings[0].StartTime,
		myBookingsPreviousExpected.Bookings[0].EndTime,
		myBookingsPreviousExpected.Bookings[0].Status,
		myBookingsPreviousExpected.Bookings[0].TotalPrice,
	).
	AddRow(
		myBookingsPreviousExpected.Bookings[1].ID,
		myBookingsPreviousExpected.Bookings[1].PlaceID,
		myBookingsPreviousExpected.Bookings[1].PlaceName,
		myBookingsPreviousExpected.Bookings[1].PlaceImage,
		myBookingsPreviousExpected.Bookings[1].Date,
		myBookingsPreviousExpected.Bookings[1].StartTime,
		myBookingsPreviousExpected.Bookings[1].EndTime,
		myBookingsPreviousExpected.Bookings[1].Status,
		myBookingsPreviousExpected.Bookings[1].TotalPrice,
	)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now()
	`)).
		WithArgs(localID).
		WillReturnRows(rows)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Equal(t, myBookingsPreviousExpected, myBookingsPreviousRetrieve)
	assert.NotNil(t, myBookingsPreviousRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetMyBookingsPreviousWithPaginationEmpty(t *testing.T) {
	myBookingsPreviousExpected := &List{
		Bookings:     make([]Booking, 0),
		TotalCount: 0,
	}
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Equal(t, myBookingsPreviousExpected, myBookingsPreviousRetrieve)
	assert.NotNil(t, myBookingsPreviousRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetMyBookingsPreviousWithPaginationEmptyWhenCount(t *testing.T) {
	myBookingsPreviousExpected := &List{
		Bookings:     make([]Booking, 0),
		TotalCount: 0,
	}
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
	NewRows([]string{"id", "place_id", "place_name", "place_image", "date", "start_time", "end_time", "status", "total_price"})
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now()
	`)).
		WillReturnError(sql.ErrNoRows)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Equal(t, myBookingsPreviousExpected, myBookingsPreviousRetrieve)
	assert.NotNil(t, myBookingsPreviousRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetMyBookingsPreviousWithPaginationError(t *testing.T) {
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Nil(t, myBookingsPreviousRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetMyBookingsPreviousWithPaginationErrorWhenCount(t *testing.T) {
	localID := "abc"
	params := BookingsListRequest{
		Limit: 10,
		Page:  1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
	NewRows([]string{"id", "place_id", "place_name", "place_image", "date", "start_time", "end_time", "status", "total_price"})
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, bookings.total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now() 
		ORDER BY bookings.end_time DESC LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.end_time <= now()
	`)).
		WillReturnError(sql.ErrConnDone)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Nil(t, myBookingsPreviousRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}
