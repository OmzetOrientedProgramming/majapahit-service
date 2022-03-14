package auth

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo PostgreSQL for auth module
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
	CheckPhoneNumber(phoneNumber string) (bool, error)

	CreateCustomer(customer Customer) (*Customer, error)

	GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error)
}

func (r repo) CheckPhoneNumber(phoneNumber string) (bool, error) {
	var phoneNumberResult string
	err := r.db.Get(&phoneNumberResult, "SELECT phone_number FROM users WHERE phone_number=$1 LIMIT 1", phoneNumber)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	}
	return false, errors.Wrap(ErrInternalServer, err.Error())
}

func (r repo) CreateCustomer(customer Customer) (*Customer, error) {
	// Inserting to parent user table
	userModel := &UserModel{
		PhoneNumber: customer.PhoneNumber,
		Name:        customer.Name,
		Status:      customer.Status,
	}
	query := `
		INSERT INTO users (phone_number, name, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var lastInsertID int
	err := r.db.QueryRowx(query, userModel.PhoneNumber, userModel.Name, userModel.Status).Scan(&lastInsertID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	// Inserting to child customer table
	customerModel := &CustomerModel{
		DateOfBirth: customer.DateOfBirth,
		Gender:      customer.Gender,
		UserID:      lastInsertID,
	}
	query = `
	INSERT INTO customers (date_of_birth, gender, user_id)
	VALUES (:date_of_birth, :gender, :user_id) 
	`
	if _, err = r.db.NamedExec(query, customerModel); err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &Customer{
		ID:          lastInsertID,
		DateOfBirth: customer.DateOfBirth,
		Gender:      customer.Gender,
		PhoneNumber: customer.PhoneNumber,
		Name:        customer.Name,
		Status:      customer.Status,
	}, nil
}

func (r repo) GetCustomerByPhoneNumber(phoneNumber string) (*Customer, error) {
	var userModel UserModel
	query := "SELECT * FROM users WHERE phone_number=$1"
	if err := r.db.Get(&userModel, query, phoneNumber); err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	var customerModel CustomerModel
	query = "SELECT * FROM customers WHERE user_id=$1"
	if err := r.db.Get(&customerModel, query, userModel.ID); err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &Customer{
		ID:          userModel.ID,
		DateOfBirth: customerModel.DateOfBirth,
		Gender:      customerModel.Gender,
		PhoneNumber: userModel.PhoneNumber,
		Name:        userModel.Name,
		Status:      1,
	}, nil
}
