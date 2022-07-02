package review

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_InsertBookingReview(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("Successfully insert review", func(t *testing.T) {
		userID := 1
		placeID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID		: userID,
			PlaceID		: placeID,
			BookingID	: bookingID,
			Content		: content,
			Rating		: rating,
		}

		query := `
		INSERT INTO reviews (user_id, place_id, booking_id, content, rating)
		VALUES (?, ?, ?, ?, ?);
		`

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
			review.UserID, review.PlaceID, review.BookingID, review.Content, review.Rating).
			WillReturnResult(driver.ResultNoRows)

		err := repoMock.InsertBookingReview(review)
		assert.Nil(t, err)
	})

	t.Run("Internal server error while inserting review", func(t *testing.T) {
		userID := 1
		placeID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID		: userID,
			PlaceID		: placeID,
			BookingID	: bookingID,
			Content		: content,
			Rating		: rating,
		}

		query := `
		INSERT INTO reviews (user_id, place_id, booking_id, content, rating)
		VALUES (?, ?, ?, ?, ?);
		`

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
			review.UserID, review.PlaceID, review.BookingID, review.Content, review.Rating).
			WillReturnError(sql.ErrTxDone)

		err := repoMock.InsertBookingReview(review)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})
}

func TestRepo_RetrievePlaceID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("Successfully retrieve place ID", func(t *testing.T) {
		bookingID := 1

		query := `
		SELECT place_id FROM bookings WHERE id=$1 LIMIT 1;
		`

		rows := mock.
			NewRows([]string{"place_id"}).
			AddRow(5)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnRows(rows)

		placeID, err := repoMock.RetrievePlaceID(1)
		assert.NoError(t, err)
		assert.Equal(t, 5, *placeID)
	})

	t.Run("No place ID to retrieve", func(t *testing.T) {
		bookingID := 1

		query := `
		SELECT place_id FROM bookings WHERE id=$1 LIMIT 1;
		`

		rows := mock.
			NewRows([]string{"place_id"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnRows(rows)

		placeID, err := repoMock.RetrievePlaceID(1)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, placeID)
	})
}

func TestRepo_CheckBookingStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("Booking is eligible to be reviewed", func(t *testing.T) {
		bookingID := 1

		query := `
		SELECT status FROM bookings WHERE id=$1 LIMIT 1;
		`

		rows := mock.
			NewRows([]string{"status"}).
			AddRow(3)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnRows(rows)

		isEligible, err := repoMock.CheckBookingStatus(bookingID)
		assert.NoError(t, err)
		assert.Equal(t, true, isEligible)
	})

	t.Run("No booking to check", func(t *testing.T) {
		bookingID := 1

		query := `
		SELECT status FROM bookings WHERE id=$1 LIMIT 1;
		`

		rows := mock.
			NewRows([]string{"status"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnRows(rows)

		isEligible, err := repoMock.CheckBookingStatus(bookingID)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
		assert.Equal(t, false, isEligible)
	})

	t.Run("Booking is not eligible to be reviewed", func(t *testing.T) {
		bookingID := 1

		query := `
		SELECT status FROM bookings WHERE id=$1 LIMIT 1;
		`

		rows := mock.
			NewRows([]string{"status"}).
			AddRow(4)
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnRows(rows)

		isEligible, err := repoMock.CheckBookingStatus(bookingID)
		assert.NoError(t, errors.Cause(err))
		assert.Equal(t, false, isEligible)
	})
}

func TestRepo_UpdateBookingStatus(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("Booking status is successfully updated", func(t *testing.T) {
		bookingID := 1

		query := `
		UPDATE bookings
		SET status = 5
		WHERE id = $1
		`

		mock.NewRows([]string{"status"}).AddRow(3)

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnResult(sqlmock.NewResult(int64(bookingID), 1))

		err := repoMock.UpdateBookingStatus(bookingID)
		assert.NoError(t, err)
	})

	t.Run("Internal server error", func(t *testing.T) {
		bookingID := 1

		query := `
		UPDATE bookings
		SET status = 5
		WHERE id = $1
		`

		mock.NewRows([]string{"status"})

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(bookingID).
			WillReturnError(ErrInternalServer)

		err := repoMock.UpdateBookingStatus(bookingID)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})
}