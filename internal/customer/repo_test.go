package customer

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

func TestRepo_PutEditCustomer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)

	t.Run("Successfully put customer", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		query := `
			UPDATE users
			SET image = $1,
					name = $2
			WHERE id = $3
		`
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.ProfilePicture, customer.Name, customer.ID).WillReturnResult(driver.ResultNoRows)

		query = `
			UPDATE customers
			SET date_of_birth = $1,
					gender = $2
			WHERE user_id = $3
		`

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.DateOfBirth, customer.Gender, customer.ID).WillReturnResult(driver.ResultNoRows)

		err := repoMock.PutEditCustomer(customer)
		assert.Nil(t, err)
	})

	t.Run("Error user with given ID not found", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		query := `
			UPDATE users
			SET image = $1,
					name = $2
			WHERE id = $3
		`
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.ProfilePicture, customer.Name, customer.ID).WillReturnError(sql.ErrNoRows)

		err := repoMock.PutEditCustomer(customer)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
	})

	t.Run("Error internal server error when updating user row", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		query := `
			UPDATE users
			SET image = $1,
					name = $2
			WHERE id = $3
		`
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.ProfilePicture, customer.Name, customer.ID).WillReturnError(sql.ErrTxDone)

		err := repoMock.PutEditCustomer(customer)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})

	t.Run("Error customer with given ID not found", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		query := `
			UPDATE users
			SET image = $1,
					name = $2
			WHERE id = $3
		`
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.ProfilePicture, customer.Name, customer.ID).WillReturnResult(driver.ResultNoRows)

		query = `
			UPDATE customers
			SET date_of_birth = $1,
					gender = $2
			WHERE user_id = $3
		`

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.DateOfBirth, customer.Gender, customer.ID).WillReturnError(sql.ErrNoRows)

		err := repoMock.PutEditCustomer(customer)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidation, errors.Cause(err))
	})

	t.Run("Error internal server error when updating customer row", func(t *testing.T) {
		userID := 1
		name := "Customer 123"
		profilePicture := "https://asset-a.grid.id//crop/0x0:0x0/700x465/photo/bobofoto/original/17064_2-cara-untuk-melakukan-sikap-kayang.JPG"
		dateOfBirthString := "2000-04-09"
		dateOfBirth, _ := time.Parse(util.DateLayout, dateOfBirthString)
		gender := 1

		customer := EditCustomerRequest{
			ID:                userID,
			Name:              name,
			ProfilePicture:    profilePicture,
			DateOfBirth:       dateOfBirth,
			DateOfBirthString: dateOfBirthString,
			Gender:            gender,
		}

		query := `
			UPDATE users
			SET image = $1,
					name = $2
			WHERE id = $3
		`
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.ProfilePicture, customer.Name, customer.ID).WillReturnResult(driver.ResultNoRows)

		query = `
			UPDATE customers
			SET date_of_birth = $1,
					gender = $2
			WHERE user_id = $3
		`

		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(customer.DateOfBirth, customer.Gender, customer.ID).WillReturnError(sql.ErrTxDone)

		err := repoMock.PutEditCustomer(customer)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})
}

func TestRepo_RetrieveCustomerProfile(t *testing.T) {
	userID := 1
	customerProfileExpected := &Profile{
		PhoneNumber:    "08123456789",
		Name:           "test_name_profile",
		Gender:         0,
		DateOfBirth:    time.Date(2001, 6, 10, 0, 0, 0, 0, time.UTC),
		ProfilePicture: "test_image_profile",
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
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
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
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
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
		assert.Equal(t, ErrInternalServer, errors.Cause(err))
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
