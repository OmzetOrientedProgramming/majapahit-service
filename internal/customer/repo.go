package customer

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
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
	RetrieveCustomerProfile(userID int) (*Profile, error)
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

// RetrieveCustomerProfile returns Profile struct for later to be the data of response body
func (r repo) RetrieveCustomerProfile(userID int) (*Profile, error) {
	var userModel auth.UserModel
	query := "SELECT phone_number, name, image FROM users WHERE id = $1"
	if err := r.db.Get(&userModel, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, "user not found")
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	var customerModel auth.CustomerModel
	query = "SELECT gender, date_of_birth FROM customers WHERE user_id = $1"
	if err := r.db.Get(&customerModel, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, "user not found")
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &Profile{
		PhoneNumber:    userModel.PhoneNumber,
		Name:           userModel.Name,
		Gender:         customerModel.Gender,
		DateOfBirth:    customerModel.DateOfBirth,
		ProfilePicture: userModel.Image,
	}, nil
}
