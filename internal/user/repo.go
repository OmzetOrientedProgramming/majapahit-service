package user

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
	GetUserIDByLocalID(localID string) (*Model, error)
}

func (r repo) GetUserIDByLocalID(localID string) (*Model, error) {
	var user Model
	err := r.db.Get(&user, "SELECT id, phone_number, name, status, COALESCE(email, '') as email, COALESCE(firebase_local_id, '') as firebase_local_id, created_at, updated_at FROM users WHERE firebase_local_id=$1", localID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(ErrNotFound, err.Error())
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &user, nil
}
