package customerbooking

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetListCustomerBookingWwithPaginationSuccess(t *testing.T) {
	listCustomerBookingExpected := &List{
		CustomerBookings: []CustomerBooking{
			{
				ID:          	1,
				Capacity:       10,
				Date:       	"test date",
				StartTime: 		"test start time",
				EndTime:		"test end time",
			},
			{
				ID:          	2,
				Capacity:       10,
				Date:       	"test date",
				StartTime: 		"test start time",
				EndTime: 		"test end time",
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
		NewRows([]string{"id", "capacity", "date", "start_time", "end_time"}).
		AddRow(listCustomerBookingExpected.CustomerBookings[0].ID,
			listCustomerBookingExpected.CustomerBookings[0].Capacity,
			listCustomerBookingExpected.CustomerBookings[0].Date,
			listCustomerBookingExpected.CustomerBookings[0].StartTime,
			listCustomerBookingExpected.CustomerBookings[0].EndTime).
		AddRow(listCustomerBookingExpected.CustomerBookings[1].ID,
			listCustomerBookingExpected.CustomerBookings[1].Capacity,
			listCustomerBookingExpected.CustomerBookings[1].Date,
			listCustomerBookingExpected.CustomerBookings[1].StartTime,
			listCustomerBookingExpected.CustomerBookings[1].EndTime)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.id, b.capacity, b.date, b.start_time, b.end_time FROM bookings b WHERE b.place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM bookings WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)

	// Test
	listItemResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}