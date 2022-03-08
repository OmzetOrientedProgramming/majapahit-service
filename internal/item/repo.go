package item

import (
	"database/sql"
	"fmt"

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
	GetItemById(place_id int, item_id int) (*Item, error)
}

func (r repo) GetListItem(place_id int, name string) (*ListItem, error) {
	var listItem ListItem
	var listQuery []interface{}
	n := 1
	listItem.Items = make([]Item, 0)

	query := "SELECT id, name, image, price, description FROM items WHERE "
	
	if name != "" {
		query += fmt.Sprintf("name LIKE $%d AND ", n)
		n += 1
		listQuery = append(listQuery, "%"+name+"%")
	}

	query += fmt.Sprintf("place_id = $%d", n)
	listQuery = append(listQuery, place_id)
	err := r.db.Select(&listItem.Items, query, listQuery...)
	
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}
	
	return &listItem, nil
}

func (r repo) GetItemById(place_id int, item_id int) (*Item, error) {
	var item Item
	item = Item{}

	query := "SELECT id, name, image, price, description FROM items WHERE place_id = $1 AND id = $2"
	err := r.db.Get(&item, query, place_id, item_id)

	if err != nil {
		if err == sql.ErrNoRows {
			item = Item{}
			return &item, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &item, nil
}