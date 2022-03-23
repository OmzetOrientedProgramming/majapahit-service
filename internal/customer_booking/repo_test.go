package customerbooking

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetListCustomerBookingWwithPaginationSuccess(t *testing.T) {
	listCustomerBookingExpected := &List{
		CustomerBookings: []CustomerBooking{
			{
				ID:           1,
				CustomerName: "test name",
				Capacity:     10,
				Date:         "test date",
				StartTime:    "test start time",
				EndTime:      "test end time",
			},
			{
				ID:           2,
				CustomerName: "test name",
				Capacity:     10,
				Date:         "test date",
				StartTime:    "test start time",
				EndTime:      "test end time",
			},
		},
		TotalCount: 10,
	}

	params := ListRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
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
		NewRows([]string{"id", "name", "capacity", "date", "start_time", "end_time"}).
		AddRow(listCustomerBookingExpected.CustomerBookings[0].ID,
			listCustomerBookingExpected.CustomerBookings[0].CustomerName,
			listCustomerBookingExpected.CustomerBookings[0].Capacity,
			listCustomerBookingExpected.CustomerBookings[0].Date,
			listCustomerBookingExpected.CustomerBookings[0].StartTime,
			listCustomerBookingExpected.CustomerBookings[0].EndTime).
		AddRow(listCustomerBookingExpected.CustomerBookings[1].ID,
			listCustomerBookingExpected.CustomerBookings[1].CustomerName,
			listCustomerBookingExpected.CustomerBookings[1].Capacity,
			listCustomerBookingExpected.CustomerBookings[1].Date,
			listCustomerBookingExpected.CustomerBookings[1].StartTime,
			listCustomerBookingExpected.CustomerBookings[1].EndTime)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4")).
		WithArgs(params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM bookings WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}

func TestRepo_GetListCustomerBookingWithPaginationError(t *testing.T) {
	params := ListRequest{
		Limit:   10,
		Page:    1,
		State:   1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4")).
		WithArgs(params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Nil(t, listCustomerBookingResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListCustomerBookingWithPaginationCountError(t *testing.T) {
	params := ListRequest{
		Limit:   10,
		Page:    1,
		State:   1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "capacity", "date", "start_time", "end_time"}).
		AddRow("1", "test name", 1, "test date", "test start time", "test end time")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4")).
		WithArgs(params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM bookings WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrConnDone)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Nil(t, listCustomerBookingResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListCustomerBookingWithPaginationEmpty(t *testing.T) {
	listCustomerBookingExpected := &List{
		CustomerBookings: make([]CustomerBooking, 0),
	}

	params := ListRequest{
		Limit:   10,
		Page:    1,
		State:   1,
		PlaceID: 1,
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4")).
		WithArgs(params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)

}

func TestRepo_GetListItemWithPaginationCountEmpty(t *testing.T) {
	listCustomerBookingExpected := &List{
		CustomerBookings: make([]CustomerBooking, 0),
		TotalCount:       0,
	}

	params := ListRequest{
		Limit:   10,
		Page:    1,
		State:   1,
		PlaceID: 1,
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
		NewRows([]string{"id", "name", "capacity", "date", "start_time", "end_time"}).
		AddRow("1", "test name", 1, "test date", "test start time", "test end time")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time FROM bookings b, users u WHERE b.place_id = $1 AND u.id = b.user_id AND b.status = $2 LIMIT $3 OFFSET $4")).
		WithArgs(params.PlaceID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM bookings WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrNoRows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}
