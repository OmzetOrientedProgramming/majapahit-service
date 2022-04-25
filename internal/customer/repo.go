package customer

import (
	"github.com/jmoiron/sqlx"
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
	return nil
}
