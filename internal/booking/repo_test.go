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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT items.name as name, items.image as image, items.qty as qty, items.price as price FROM items INNER JOIN booking_items ON items.id = booking_items.item_id WHERE booking_items.booking_id = $1")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT items.name as name, items.image as image, items.qty as qty, items.price as price FROM items INNER JOIN booking_items ON items.id = booking_items.item_id WHERE booking_items.booking_id = $1")).
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
