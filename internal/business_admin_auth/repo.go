package businessadminauth

import (
	"github.com/jmoiron/sqlx"
)

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type repo struct {
	db *sqlx.DB
}

type Repo interface {
}
