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
	GetListItem(place_id int, name string) (*ListItem, error)
	GetItemById(item_id int) (*Item, error)
}

func (r repo) GetListItem(place_id int, name string) (*ListItem, error) {
	var listItem ListItem
	listItem.Items = make([]Item, 0)

	query := "SELECT id, name, image, price, description FROM items WHERE "
	
	if name != "" {
		query = query + "name LIKE '%$1% AND "
	}

	query = query + "place_id = $2"
	err := r.db.Select(&listItem.Items, query, name, place_id)
	
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}
	
	return &listItem, nil
}
func (r repo) GetItemById(item_id int) (*Item, error) {
	var item Item
	item = Item{}

	query := "SELECT id, name, image, price, description FROM items WHERE item_id = $1"
	err := r.db.Get(&item, query, item_id)

	if err != nil {
		if err == sql.ErrNoRows {
			item = Item{}
			return &item, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &item, nil
}