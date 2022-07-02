package auth

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
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
		Gender:      "undefined",
		PhoneNumber: "081223902345",
		Name:        "John Doe",
		Status:      "customer",
	}

	mockCustomerMale := Customer{
		DateOfBirth: time.Time{},
		Gender:      "male",
		PhoneNumber: "081223902345",
		Name:        "John Doe",
		Status:      "business admin",
	}

	mockCustomerFemale := Customer{
		DateOfBirth: time.Time{},
		Gender:      "female",
		PhoneNumber: "081223902345",
		Name:        "John Doe",
		Status:      "customer",
	}

	t.Run("success male", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomerMale.PhoneNumber, mockCustomerMale.Name, 1, mockCustomerMale.LocalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomerMale.PhoneNumber,
			Name:        mockCustomerMale.Name,
			Status:      mockCustomerMale.Status,
			Gender:      mockCustomerMale.Gender,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdCustomer)
		assert.Equal(t, createdCustomer, &Customer{
			ID:          1,
			DateOfBirth: time.Time{},
			Gender:      "male",
			PhoneNumber: "081223902345",
			Name:        "John Doe",
			Status:      "business admin",
		})
	})

	t.Run("success female", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomerFemale.PhoneNumber, mockCustomerFemale.Name, 0, mockCustomerFemale.LocalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomerFemale.PhoneNumber,
			Name:        mockCustomerFemale.Name,
			Status:      mockCustomerFemale.Status,
			Gender:      mockCustomerFemale.Gender,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdCustomer)
		assert.Equal(t, createdCustomer, &Customer{
			ID:          1,
			DateOfBirth: time.Time{},
			Gender:      "female",
			PhoneNumber: "081223902345",
			Name:        "John Doe",
			Status:      "customer",
		})
	})

	t.Run("success", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, 0, mockCustomer.LocalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdCustomer, err := repoMock.CreateCustomer(Customer{
			PhoneNumber: mockCustomer.PhoneNumber,
			Name:        mockCustomer.Name,
			Status:      mockCustomer.Status,
			Gender:      mockCustomer.Gender,
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdCustomer)
		assert.Equal(t, createdCustomer, &Customer{
			ID:          1,
			DateOfBirth: time.Time{},
			Gender:      "undefined",
			PhoneNumber: "081223902345",
			Name:        "John Doe",
			Status:      "customer",
		})
	})

	t.Run("error inserting to user table", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, mockCustomer.Status, mockCustomer.LocalID).
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

func TestRepo_CreateCustomerErrorOnInsertToCustomerTable(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mockCustomer := Customer{
		DateOfBirth: time.Time{},
		Gender:      "undefined",
		PhoneNumber: "081223902345",
		Name:        "John Doe",
		Status:      "customer",
	}

	t.Run("error inserting to customer table", func(t *testing.T) {
		query := `
		INSERT INTO users (.+)
		VALUES (.+)
		RETURNING id
		`
		mock.ExpectQuery(query).
			WithArgs(mockCustomer.PhoneNumber, mockCustomer.Name, 0, mockCustomer.LocalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		query = `
		INSERT INTO customers (.+)
		VALUES (.+) 
		`
		mock.ExpectExec(query).WithArgs(mockCustomer.DateOfBirth, mockCustomer.Gender, util.StatusCustomer).
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

func TestRepo_GetCustomerByPhoneNumber(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("error getting customer data", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status"}).
			AddRow(1, "081223901234", "Bambang", 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE phone_number=$1")).
			WithArgs("081223901234").
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM customers WHERE user_id=$1")).
			WithArgs(1).
			WillReturnError(ErrInternalServer)

		customer, err := repoMock.GetCustomerByPhoneNumber("081223901234")
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, customer)
	})

	t.Run("success male", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status"}).
			AddRow(1, "081223901234", "Bambang", 0)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE phone_number=$1")).
			WithArgs("081223901234").
			WillReturnRows(rows)

		rows = mock.
			NewRows([]string{"id", "date_of_birth", "gender", "user_id"}).
			AddRow(1, time.Time{}, 1, 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM customers WHERE user_id=$1")).
			WithArgs(1).
			WillReturnRows(rows)

		customer, err := repoMock.GetCustomerByPhoneNumber("081223901234")
		assert.NoError(t, err)
		assert.NotNil(t, customer)
	})

	t.Run("success undefined gender", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status"}).
			AddRow(1, "081223901234", "Bambang", 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE phone_number=$1")).
			WithArgs("081223901234").
			WillReturnRows(rows)

		rows = mock.
			NewRows([]string{"id", "date_of_birth", "gender", "user_id"}).
			AddRow(1, time.Time{}, 0, 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM customers WHERE user_id=$1")).
			WithArgs(1).
			WillReturnRows(rows)

		customer, err := repoMock.GetCustomerByPhoneNumber("081223901234")
		assert.NoError(t, err)
		assert.NotNil(t, customer)
	})

	t.Run("success female gender", func(t *testing.T) {
		rows := mock.
			NewRows([]string{"id", "phone_number", "name", "status"}).
			AddRow(1, "081223901234", "Bambang", 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE phone_number=$1")).
			WithArgs("081223901234").
			WillReturnRows(rows)

		rows = mock.
			NewRows([]string{"id", "date_of_birth", "gender", "user_id"}).
			AddRow(1, time.Time{}, 2, 1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM customers WHERE user_id=$1")).
			WithArgs(1).
			WillReturnRows(rows)

		customer, err := repoMock.GetCustomerByPhoneNumber("081223901234")
		assert.NoError(t, err)
		assert.NotNil(t, customer)
	})

	t.Run("error getting user data", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, name, status FROM users WHERE phone_number=$1 LIMIT 1")).
			WithArgs("081223902345").
			WillReturnError(ErrInternalServer)

		customer, err := repoMock.GetCustomerByPhoneNumber("081223902345")
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
		assert.Nil(t, customer)
	})
}
