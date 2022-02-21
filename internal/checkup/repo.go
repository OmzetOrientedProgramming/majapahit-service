package checkup

import (
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

type Repo interface {
	GetApplicationCheckUp() (bool, error)
}

func (r repo) GetApplicationCheckUp() (bool, error) {
	if r.db == nil {
		return false, errors.Wrap(ErrPostgreSQLNotConnected, "postgreSQL not connected")
	}

	err := r.db.Ping()
	if err != nil {
		return false, errors.Wrap(ErrPingDBFailed, err.Error())
	}

	return true, nil
}
