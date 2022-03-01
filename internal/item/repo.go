package item

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	GetListItem(place_id int) (*ListItem, error)
	GetItem() (*Item, error)
}

func (r repo) GetListItem(place_id int) (*ListItem, error) {
	var listItem ListItem
	listItem.Items = make([]Item, 0)

	query := "SELECT id, name, price, description FROM items WHERE place_id = $1"
	err := r.db.Select(&listItem.Items, query, place_id)
	
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}
	
	return &listItem, nil
}

func (r repo) GetItem() (*Item, error) {
	var item Item
	// Do something here
	return &item, nil
}