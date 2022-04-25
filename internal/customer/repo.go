package customer

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
)

// NewRepo used to initialize repo
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

// Repo interface for defining function that must have by repo
type Repo interface {
	RetrieveCustomerProfile(userID int) (*Profile, error)
}

// RetrieveCustomerProfile returns Profile struct for later to be the data of response body
func (r repo) RetrieveCustomerProfile(userID int) (*Profile, error) {
	var userModel auth.UserModel
	query := "SELECT phone_number, name, image FROM users WHERE id = $1"
	if err := r.db.Get(&userModel, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, "user not found")
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	var customerModel auth.CustomerModel
	query = "SELECT gender, date_of_birth FROM customers WHERE user_id = $1"
	if err := r.db.Get(&customerModel, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(ErrNotFound, "user not found")
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &Profile{
		PhoneNumber:    userModel.PhoneNumber,
		Name:           userModel.Name,
		Gender:         customerModel.Gender,
		DateOfBirth:    customerModel.DateOfBirth,
		ProfilePicture: userModel.Image,
	}, nil
}