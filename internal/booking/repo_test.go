package booking

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

func TestRepo_GetListCustomerBookingWwithPaginationSuccess(t *testing.T) {
	listCustomerBookingExpected := &ListBooking{
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
		Limit:  10,
		Page:   1,
		UserID: 1,
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
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
		FROM bookings b, users u, places p 
		WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
		ORDER BY b.date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1 AND b.status = $2")).
		WithArgs(params.UserID, params.State).
		WillReturnRows(rows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}

func TestRepo_GetListCustomerBookingWithPaginationError(t *testing.T) {
	params := ListRequest{
		Limit:  10,
		Page:   1,
		State:  1,
		UserID: 1,
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
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
		FROM bookings b, users u, places p 
		WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
		ORDER BY b.date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Nil(t, listCustomerBookingResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListCustomerBookingWithPaginationCountError(t *testing.T) {
	params := ListRequest{
		Limit:  10,
		Page:   1,
		State:  1,
		UserID: 1,
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
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
		FROM bookings b, users u, places p 
		WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
		ORDER BY b.date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1 AND b.status = $2")).
		WithArgs(params.UserID, params.State).
		WillReturnError(sql.ErrConnDone)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Nil(t, listCustomerBookingResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListCustomerBookingWithPaginationEmpty(t *testing.T) {
	listCustomerBookingExpected := &ListBooking{
		CustomerBookings: make([]CustomerBooking, 0),
	}

	params := ListRequest{
		Limit:  10,
		Page:   1,
		State:  1,
		UserID: 1,
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
		SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
		FROM bookings b, users u, places p 
		WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
		ORDER BY b.date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)

}

func TestRepo_GetListItemWithPaginationCountEmpty(t *testing.T) {
	listCustomerBookingExpected := &ListBooking{
		CustomerBookings: make([]CustomerBooking, 0),
		TotalCount:       0,
	}

	params := ListRequest{
		Limit:  10,
		Page:   1,
		State:  1,
		UserID: 1,
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
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT b.id, u.name, b.capacity, b.date, b.start_time, b.end_time 
		FROM bookings b, users u, places p 
		WHERE p.user_id = $1 AND p.id = b.place_id AND u.id = b.user_id AND b.status = $2 
		ORDER BY b.date DESC LIMIT $3 OFFSET $4`)).
		WithArgs(params.UserID, params.State, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1 AND b.status = $2")).
		WithArgs(params.UserID, params.State).
		WillReturnError(sql.ErrNoRows)

	// Test
	listCustomerBookingResult, err := repoMock.GetListCustomerBookingWithPagination(params)
	assert.Equal(t, listCustomerBookingExpected, listCustomerBookingResult)
	assert.NotNil(t, listCustomerBookingResult)
	assert.NoError(t, err)
}

func TestRepo_GetBookingData(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("success", func(t *testing.T) {
		params := GetBookingDataParams{
			PlaceID:   0,
			StartDate: time.Time{},
			EndDate:   time.Time{},
			StartTime: time.Time{},
		}

		query := `SELECT id, date, start_time, end_time, capacity 
				FROM bookings 
				WHERE place_id = $1
				AND (status = $2 or status = $3)
				AND date >= $4 
				AND date <= $5`

		rows := mock.
			NewRows([]string{"id", "date", "start_time", "end_time", "capacity"}).
			AddRow(1, time.Now(), time.Now(), time.Now(), 10)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(params.PlaceID, util.BookingBelumMembayar, util.BookingBerhasil, params.StartDate, params.EndDate).
			WillReturnRows(rows)

		bookingData, err := repoMock.GetBookingData(params)
		assert.NotNil(t, bookingData)
		assert.Nil(t, err)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		params := GetBookingDataParams{
			PlaceID:   0,
			StartDate: time.Time{},
			EndDate:   time.Time{},
			StartTime: time.Time{},
		}

		query := `SELECT id, date, start_time, end_time, capacity 
				FROM bookings 
				WHERE place_id = $1
				AND status= $2
				AND date >= $3 
				AND date <= $4`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(params.PlaceID, util.BookingBelumMembayar, params.StartDate, params.EndDate).
			WillReturnError(ErrInternalServerError)

		bookingData, err := repoMock.GetBookingData(params)
		assert.Nil(t, bookingData)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestRepo_GetTimeSlotsData(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("success", func(t *testing.T) {
		placeID := 1
		selectedDate := time.Now()

		query := `SELECT id, start_time, end_time, day
				FROM time_slots 
				WHERE place_id = $1 
				AND (day = $2)
				ORDER BY day, start_time`

		rows := mock.
			NewRows([]string{"id", "start_time", "end_time", "day"}).
			AddRow(1, time.Now(), time.Now(), 1)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(placeID, int(selectedDate.Weekday())).
			WillReturnRows(rows)

		timeSlot, err := repoMock.GetTimeSlotsData(placeID, selectedDate)
		assert.NotNil(t, timeSlot)
		assert.Nil(t, err)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		placeID := 1
		selectedDate := time.Now()

		query := `SELECT id, start_time, end_time, day
				FROM time_slots 
				WHERE place_id = $1 
				AND day = $2
				ORDER BY start_time`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(placeID, int(selectedDate.Weekday())).
			WillReturnError(ErrInternalServerError)

		timeSlotsData, err := repoMock.GetTimeSlotsData(placeID, selectedDate)
		assert.Nil(t, timeSlotsData)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestRepo_GetPlaceCapacity(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("success", func(t *testing.T) {
		placeID := 1
		placeCapacity := 10
		openHour, _ := time.Parse(util.TimeLayout, "08:00:00")

		query := `SELECT capacity, open_hour FROM places WHERE id = $1`

		rows := mock.
			NewRows([]string{"capacity", "open_hour"}).
			AddRow(placeCapacity, openHour)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(placeID).
			WillReturnRows(rows)

		placeCapacityRes, err := repoMock.GetPlaceCapacity(placeID)
		assert.Equal(t, &PlaceOpenHourAndCapacity{
			OpenHour: openHour,
			Capacity: placeCapacity,
		}, placeCapacityRes)
		assert.Nil(t, err)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		placeID := 1

		query := `SELECT capacity FROM places WHERE id = $1`

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(placeID).
			WillReturnError(ErrInternalServerError)

		_, err := repoMock.GetPlaceCapacity(placeID)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestRepo_CheckedItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := []CheckedItemParams{
			{
				ID:      1,
				PlaceID: 1,
			},
			{
				ID:      2,
				PlaceID: 1,
			},
		}
		rows := mock.NewRows([]string{"id", "place_id"}).AddRow("1", "1").AddRow("2", "1")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, place_id FROM items WHERE place_id = $1 AND (id = $2 OR id = $3)")).WithArgs(1, 1, 2).WillReturnRows(rows)

		item, isMatch, err := repo.CheckedItem(input)
		assert.Nil(t, err)
		assert.True(t, isMatch)
		assert.Equal(t, &input, item)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := []CheckedItemParams{
			{
				ID:      1,
				PlaceID: 1,
			},
			{
				ID:      2,
				PlaceID: 1,
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, place_id FROM items WHERE place_id = $1 AND (id = $2 OR id = $3)")).WithArgs(1, 1, 2).WillReturnError(ErrInternalServerError)

		item, isMatch, err := repo.CheckedItem(input)
		assert.NotNil(t, err)
		assert.Nil(t, item)
		assert.False(t, isMatch)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed item not match", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := []CheckedItemParams{
			{
				ID:      1,
				PlaceID: 1,
			},
			{
				ID:      2,
				PlaceID: 1,
			},
		}

		rows := mock.NewRows([]string{"id", "place_id"}).AddRow("1", "1")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, place_id FROM items WHERE place_id = $1 AND (id = $2 OR id = $3)")).WithArgs(1, 1, 2).WillReturnRows(rows)

		item, isMatch, err := repo.CheckedItem(input)
		assert.NotNil(t, err)
		assert.NotNil(t, item)
		assert.False(t, isMatch)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestRepo_CreateBooking(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		date, _ := time.Parse(util.DateLayout, "2022-01-01")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		endTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		booking := CreateBookingParams{
			UserID:     1,
			PlaceID:    1,
			Date:       date,
			StartTime:  startTime,
			EndTime:    endTime,
			Capacity:   10,
			Status:     1,
			TotalPrice: 10000,
		}

		rows := mock.NewRows([]string{"id"}).AddRow("1")
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO 
					bookings (user_id, place_id, date, start_time, end_time, capacity, status, total_price)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id`)).
			WithArgs(booking.UserID, booking.PlaceID, booking.Date, booking.StartTime, booking.EndTime, booking.Capacity, booking.Status, booking.TotalPrice).
			WillReturnRows(rows)

		res, err := repo.CreateBooking(booking)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Equal(t, 1, res.ID)
	})

	t.Run("failed", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		date, _ := time.Parse(util.DateLayout, "2022-01-01")
		startTime, _ := time.Parse(util.TimeLayout, "08:00:00")
		endTime, _ := time.Parse(util.TimeLayout, "09:00:00")
		booking := CreateBookingParams{
			UserID:     1,
			PlaceID:    1,
			Date:       date,
			StartTime:  startTime,
			EndTime:    endTime,
			Capacity:   10,
			Status:     1,
			TotalPrice: 10000,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO 
					bookings (user_id, place_id, date, start_time, end_time, capacity, status, total_price)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id`)).
			WithArgs(booking.UserID, booking.PlaceID, booking.Date, booking.StartTime, booking.EndTime, booking.Capacity, booking.Status, booking.TotalPrice).
			WillReturnError(ErrInternalServerError)

		res, err := repo.CreateBooking(booking)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestRepo_CreateBookingItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := []CreateBookingItemsParams{
			{
				BookingID:  3,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  4,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  5,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
		}

		query := `INSERT INTO 
					booking_items (item_id, booking_id, qty, total_price)
				VALUES ($1, $2, $3, $4) , ($5, $6, $7, $8) , ($9, $10, $11, $12) `
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(input[0].ItemID, input[0].BookingID, input[0].Qty, input[0].TotalPrice,
				input[1].ItemID, input[1].BookingID, input[1].Qty, input[1].TotalPrice,
				input[2].ItemID, input[2].BookingID, input[2].Qty, input[2].TotalPrice).
			WillReturnResult(driver.ResultNoRows)

		res, err := repo.CreateBookingItems(input)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Equal(t, float64(60000), res.TotalPrice)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := []CreateBookingItemsParams{
			{
				BookingID:  3,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  4,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
			{
				BookingID:  5,
				ItemID:     1,
				TotalPrice: 20000,
				Qty:        2,
			},
		}

		query := `INSERT INTO 
					booking_items (item_id, booking_id, qty, total_price)
				VALUES ($1, $2, $3, $4) , ($5, $6, $7, $8) , ($9, $10, $11, $12) `
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(input[0].ItemID, input[0].BookingID, input[0].Qty, input[0].TotalPrice,
				input[1].ItemID, input[1].BookingID, input[1].Qty, input[1].TotalPrice,
				input[2].ItemID, input[2].BookingID, input[2].Qty, input[2].TotalPrice).
			WillReturnError(ErrInternalServerError)

		res, err := repo.CreateBookingItems(input)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestRepo_UpdateTotalPrice(t *testing.T) {
	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := UpdateTotalPriceParams{
			BookingID:  1,
			TotalPrice: 100000,
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE bookings SET total_price = $1, updated_at = NOW() WHERE id = $2")).
			WithArgs(input.TotalPrice, input.BookingID).
			WillReturnError(ErrInternalServerError)

		res, err := repo.UpdateTotalPrice(input)
		assert.False(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		repo := NewRepo(sqlxDB)

		input := UpdateTotalPriceParams{
			BookingID:  1,
			TotalPrice: 100000,
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE bookings SET total_price = $1, updated_at = NOW() WHERE id = $2")).
			WithArgs(input.TotalPrice, input.BookingID).
			WillReturnResult(driver.ResultNoRows)

		res, err := repo.UpdateTotalPrice(input)
		assert.True(t, res)
		assert.Nil(t, err)
	})
}

func TestRepo_GetDetailSuccess(t *testing.T) {
	bookingID := 1
	createdAtRow := time.Date(2021, time.Month(10), 26, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	bookingDetailExpected := &Detail{
		ID:             1,
		CustomerName:   "test nama",
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
		NewRows([]string{"id", "name", "date", "start_time", "end_time", "capacity", "status", "total_price", "created_at"}).
		AddRow(
			bookingDetailExpected.ID,
			bookingDetailExpected.CustomerName,
			bookingDetailExpected.Date,
			bookingDetailExpected.StartTime,
			bookingDetailExpected.EndTime,
			bookingDetailExpected.Capacity,
			bookingDetailExpected.Status,
			bookingDetailExpected.TotalPriceItem,
			bookingDetailExpected.CreatedAt,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, u.name, b.date, b.start_time, b.end_time, b.capacity, b.status, b.total_price, b.created_at
									   FROM bookings b, users u
									   WHERE b.id = $1 AND b.user_id = u.id`)).
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.id, u.name, b.date, b.start_time, b.end_time, b.capacity, b.status, b.total_price, b.created_at
		FROM bookings b, users u
		WHERE b.id = $1 AND b.user_id = u.id`)).
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

func TestRepo_UpdateBookingStatusSuccess(t *testing.T) {
	bookingID := 1
	newStatus := 2

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE bookings SET status = $2 WHERE id= $1")).
		WithArgs(bookingID, newStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repoMock.UpdateBookingStatus(bookingID, newStatus)
	assert.Nil(t, err)
}

func TestRepo_UpdateBookingStatusInternalServerError(t *testing.T) {
	bookingID := 1
	newStatus := 2

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE bookings SET status = $2 WHERE id= $1")).
		WithArgs(bookingID, newStatus).
		WillReturnError(sql.ErrTxDone)

	err = repoMock.UpdateBookingStatus(bookingID, newStatus)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status <= 2
		ORDER BY bookings.date asc, bookings.start_time asc
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status <= 2
		ORDER BY bookings.date asc, bookings.start_time asc
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status <= 2
		ORDER BY bookings.date asc, bookings.start_time asc
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
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
		Bookings:   make([]Booking, 0),
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
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
		Bookings:   make([]Booking, 0),
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
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
		SELECT bookings.id, bookings.place_id, places.name as place_name, places.image as place_image, bookings.date, bookings.start_time, bookings.end_time, bookings.status, places.booking_price + bookings.total_price + 3000 as total_price 
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
		ORDER BY bookings.date desc, bookings.end_time desc LIMIT $2 OFFSET $3
	`)).
		WithArgs(localID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(bookings.id)
		FROM users 
			JOIN bookings ON users.id = bookings.user_id 	
			JOIN places ON bookings.place_id = places.id 
		WHERE users.firebase_local_id = $1 AND bookings.status > 2
	`)).
		WillReturnError(sql.ErrConnDone)

	// Test
	myBookingsPreviousRetrieve, err := repoMock.GetMyBookingsPreviousWithPagination(localID, params)
	assert.Nil(t, myBookingsPreviousRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_UpdateBookingStatusByXenditID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Expectation
		repoMock := NewRepo(sqlxDB)

		query := "UPDATE bookings SET status = $2 WHERE xendit_id= $1"
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs("1", 1).WillReturnResult(driver.ResultNoRows)

		err = repoMock.UpdateBookingStatusByXenditID("1", 1)
		assert.Nil(t, err)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Expectation
		repoMock := NewRepo(sqlxDB)

		query := "UPDATE bookings SET status = $2 WHERE xendit_id= $1"
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs("1", 1).WillReturnError(ErrInternalServerError)

		err = repoMock.UpdateBookingStatusByXenditID("1", 1)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed no rows", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Expectation
		repoMock := NewRepo(sqlxDB)

		query := "UPDATE bookings SET status = $2 WHERE xendit_id= $1"
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs("1", 1).WillReturnError(sql.ErrNoRows)

		err = repoMock.UpdateBookingStatusByXenditID("1", 1)
		assert.NotNil(t, err)
		assert.Equal(t, ErrNotFound, errors.Cause(err))
	})
}

func TestRepo_InsertXenditInformation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Expectation
		params := XenditInformation{
			XenditID:    "1",
			InvoicesURL: "test.com",
			BookingID:   1,
		}
		repoMock := NewRepo(sqlxDB)

		query := "UPDATE bookings SET xendit_id = $1, invoices_url = $2 WHERE id = $3"

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(params.XenditID, params.InvoicesURL, params.BookingID).
			WillReturnResult(driver.ResultNoRows)

		isOk, err := repoMock.InsertXenditInformation(params)
		assert.Nil(t, err)
		assert.True(t, isOk)
	})

	t.Run("error internal server", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

		// Expectation
		params := XenditInformation{
			XenditID:    "1",
			InvoicesURL: "test.com",
			BookingID:   1,
		}
		repoMock := NewRepo(sqlxDB)

		query := "UPDATE bookings SET xendit_id = $1, invoices_url = $2 WHERE id = $3"

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(params.XenditID, params.InvoicesURL, params.BookingID).
			WillReturnError(ErrInternalServerError)

		isOk, err := repoMock.InsertXenditInformation(params)
		assert.NotNil(t, err)
		assert.False(t, isOk)
	})
}

func TestRepo_GetPlaceBookingPrice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		newRepo := NewRepo(sqlxDB)

		query := `SELECT COALESCE (booking_price, 0) FROM places WHERE id  = $1`

		rows := mock.NewRows([]string{"booking_price"})
		rows.AddRow(10000.0)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnRows(rows)

		resp, err := newRepo.GetPlaceBookingPrice(1)
		assert.Nil(t, err)
		assert.Equal(t, 10000.0, resp)
	})

	t.Run("failed no sql rows", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		newRepo := NewRepo(sqlxDB)

		query := `SELECT COALESCE (booking_price, 0) FROM places WHERE id  = $1`

		rows := mock.NewRows([]string{"booking_price"})
		rows.AddRow(10000.0)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		resp, err := newRepo.GetPlaceBookingPrice(1)
		assert.Equal(t, ErrNotFound, errors.Cause(err))
		assert.Equal(t, 0.0, resp)
	})

	t.Run("failed internal server error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		newRepo := NewRepo(sqlxDB)

		query := `SELECT COALESCE (booking_price, 0) FROM places WHERE id  = $1`

		rows := mock.NewRows([]string{"booking_price"})
		rows.AddRow(10000.0)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnError(ErrInternalServerError)

		resp, err := newRepo.GetPlaceBookingPrice(1)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.Equal(t, 0.0, resp)
	})
}
