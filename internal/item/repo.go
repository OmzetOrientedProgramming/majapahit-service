package item

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo used to initialize repo
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
	GetListItemWithPagination(params ListItemRequest) (*ListItem, error)
	GetItemByID(placeID int, itemID int) (*Item, error)
}

func (r repo) GetListItemWithPagination(params ListItemRequest) (*ListItem, error) {
	var listItem ListItem
	listItem.Items = make([]Item, 0)
	listItem.TotalCount = 0
	var listQuery []interface{}
	n := 1

	mainQuery := "FROM items WHERE "
	query1 := "SELECT id, name, image, price, description "
	query2 := "SELECT COUNT(id) "

	if params.Name != "" {
		mainQuery += fmt.Sprintf("name LIKE $%d AND ", n)
		n++
		listQuery = append(listQuery, "%"+params.Name+"%")
	}

	mainQuery += fmt.Sprintf("place_id = $%d LIMIT $%d OFFSET $%d", n, n+1, n+2)
	listQuery = append(listQuery, params.PlaceID, params.Limit, (params.Page-1)*params.Limit)
	err := r.db.Select(&listItem.Items, query1+mainQuery, listQuery...)

	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	err = r.db.Get(&listItem.TotalCount, query2+mainQuery, listQuery...)
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &listItem, nil
}

func (r repo) GetItemByID(placeID int, itemID int) (*Item, error) {
	var item Item
	item = Item{}

	query := "SELECT id, name, image, price, description FROM items WHERE place_id = $1 AND id = $2"
	err := r.db.Get(&item, query, placeID, itemID)

	if err != nil {
		if err == sql.ErrNoRows {
			item = Item{}
			return &item, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &item, nil
}
