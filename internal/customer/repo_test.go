package customer

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestRepo_RetrieveCustomerProfile(t *testing.T) {
	userID := 1
	customerProfileExpected := &Profile{
		PhoneNumber: 		"08123456789",
		Name:               "test_name_profile",
		Gender: 			0,
		DateOfBirth: 		time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
		ProfilePicture:     "test_image_profile",
	}

	t.Run("Both user and customer don't exist on database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.Nil(t, customerProfileActual)
	})

	t.Run("User doesn't exist on database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

		rows := mock.
		NewRows([]string{"gender", "date_of_birth"}).
		AddRow(customerProfileExpected.Gender, customerProfileExpected.DateOfBirth)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.Nil(t, customerProfileActual)
	})

	t.Run("Customer doesn't exist on database", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB) 

		rows := mock.
		NewRows([]string{"phone_number", "name", "image"}).
		AddRow(customerProfileExpected.PhoneNumber, customerProfileExpected.Name, customerProfileExpected.ProfilePicture)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
		assert.Nil(t, customerProfileActual)
	})

	t.Run("User not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB) 

		rows := mock.
		NewRows([]string{"phone_number", "name", "image"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)
		
		rows = mock.
		NewRows([]string{"gender", "date_of_birth"}).
		AddRow(customerProfileExpected.Gender, customerProfileExpected.DateOfBirth)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, ErrNotFound, errors.Cause(err))
		assert.Nil(t, customerProfileActual)
	})

	t.Run("Customer not found", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB)

		rows := mock.
		NewRows([]string{"phone_number", "name", "image"}).
		AddRow(customerProfileExpected.PhoneNumber, customerProfileExpected.Name, customerProfileExpected.ProfilePicture)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)
		
		rows = mock.
		NewRows([]string{"gender", "date_of_birth"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, ErrNotFound, errors.Cause(err))
		assert.Nil(t, customerProfileActual)
	})

	t.Run("Customer profile retrieval success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer mockDB.Close()
		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	
		repoMock := NewRepo(sqlxDB)

		rows := mock.
		NewRows([]string{"phone_number", "name", "image"}).
		AddRow(customerProfileExpected.PhoneNumber, customerProfileExpected.Name, customerProfileExpected.ProfilePicture)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number, name, image FROM users WHERE id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)
		
		rows = mock.
		NewRows([]string{"gender", "date_of_birth"}).
		AddRow(customerProfileExpected.Gender, customerProfileExpected.DateOfBirth)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT gender, date_of_birth FROM customers WHERE user_id = $1")).
			WithArgs(userID).
			WillReturnRows(rows)

		// Test
		customerProfileActual, err := repoMock.RetrieveCustomerProfile(userID)
		assert.Equal(t, customerProfileExpected, customerProfileActual)
		assert.NotNil(t, customerProfileActual)
		assert.NoError(t, err)
	})
}