package checkup

import "github.com/jmoiron/sqlx"

// NewRepo PostgreSQL for checkup module
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

type Repo interface {
	GetApplicationCheckUp() (bool, error)
}

func (r repo) GetApplicationCheckUp() (bool, error) {
	return true, nil
}
