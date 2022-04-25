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
