package item

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
	GetListItem() (*ListItem, error)
	GetItem() (*Item, error)
}

func (r repo) GetListItem() (*ListItem, error) {
	var listItem ListItem
	// Do something here
	return &listItem, nil
}

func (r repo) GetItem() (*Item, error) {
	var item Item
	// Do something here
	return &item, nil
}