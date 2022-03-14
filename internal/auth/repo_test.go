package auth

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

func TestRepo_CheckPhoneNumber(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repoMock := NewRepo(sqlxDB)

	t.Run("number exist on database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"phone_number"}).
			AddRow("081223901234")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081223901234").
			WillReturnRows(rows)

		exist, err := repoMock.CheckPhoneNumber("081223901234")
		assert.NoError(t, err)
		assert.True(t, exist)
	})

	t.Run("number doesn't exist on database", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"phone_number"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081223902345").
			WillReturnRows(rows)

		exist, err := repoMock.CheckPhoneNumber("081223902345")
		assert.NoError(t, err)
		assert.False(t, exist)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081223902345").
			WillReturnError(sql.ErrTxDone)

		exist, err := repoMock.CheckPhoneNumber("081223902345")
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.False(t, exist)
	})

}

func TestRepo_CreateCustomer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mockCustomer := Customer{
		DateOfBirth: time.Time{},
		Gender:      false,
		PhoneNumber: "081223902345",
		Name:        "John Doe",
		Status:      1,
	}

	t.Run("success", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, mockCustomer.Status).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(mockCustomer.DateOfBirth, mockCustomer.Gender, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomer.PhoneNumber,
			Name:        mockCustomer.Name,
			Status:      mockCustomer.Status,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdCustomer)
		assert.Equal(t, createdCustomer, &Customer{
			ID:          1,
			DateOfBirth: time.Time{},
			Gender:      false,
			PhoneNumber: "081223902345",
			Name:        "John Doe",
			Status:      1,
		})
	})

	t.Run("error inserting to user table", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, mockCustomer.Status).
			WillReturnError(ErrInternalServer)

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomer.PhoneNumber,
			Name:        mockCustomer.Name,
			Status:      mockCustomer.Status,
		})
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, createdCustomer)
	})

	t.Run("error inserting to customer table", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, mockCustomer.Status).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(mockCustomer.DateOfBirth, mockCustomer.Gender, 1).
			WillReturnError(ErrInternalServer)

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomer.PhoneNumber,
			Name:        mockCustomer.Name,
			Status:      mockCustomer.Status,
		})
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, createdCustomer)
	})
}
