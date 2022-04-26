package customer

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo PostgreSQL for checkup module
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

// Repo will contain all the function that can be used by repo
type Repo interface {
	PutEditCustomer(customer EditCustomerRequest) error
}

func (r repo) PutEditCustomer(customer EditCustomerRequest) error {
	// Updating user's profile picture image with corresponding ID
  query := `
    UPDATE users
    SET image = $1,
        name = $2
    WHERE id = $3
  `
  _, err := r.db.Exec(query, customer.ProfilePicture, customer.Name, customer.ID)
  if err != nil {
    if err == sql.ErrNoRows {
      return errors.Wrap(ErrInputValidation, "User ID tidak ditemukan")
    }
    return errors.Wrap(ErrInternalServer, err.Error())
  }

  query = `
    UPDATE customers
    SET date_of_birth = $1,
        gender = $2
    WHERE user_id = $3
  `

  _, err = r.db.Exec(query, customer.DateOfBirth, customer.Gender, customer.ID)
  if err != nil {
    if err == sql.ErrNoRows {
      return errors.Wrap(ErrInputValidation, "User ID tidak ditemukan")
    }
    return errors.Wrap(ErrInternalServer, err.Error())
  }

  return nil

}
