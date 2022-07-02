package item

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	GetListItemAdminWithPagination(params ListItemRequest) (*ListItem, error)
	DeleteItemAdminByID(itemID int) error
	UpdateItem(ID int, item Item) error
	CreateItem(userID int, item Item) error
}

func (r repo) GetListItemWithPagination(params ListItemRequest) (*ListItem, error) {
	var listItem ListItem
	listItem.Items = make([]Item, 0)
	listItem.TotalCount = 0
	listItem.PlaceInfo = make([]PlaceInfo, 0)
	var listQuery []interface{}
	n := 1

	mainQuery := "FROM items WHERE "
	query1 := "SELECT id, name, image, price, description "
	query2 := "SELECT COUNT(id) "
	query3 := "SELECT name, image FROM places WHERE id = $1"

	if params.Name != "" {
		mainQuery += fmt.Sprintf("LOWER(name) LIKE LOWER($%d) AND ", n)
		n++
		listQuery = append(listQuery, "%"+params.Name+"%")
	}

	tempQuery2 := mainQuery + fmt.Sprintf("place_id = $%d", n)
	mainQuery += fmt.Sprintf("place_id = $%d LIMIT $%d OFFSET $%d", n, n+1, n+2)
	listQuery = append(listQuery, params.PlaceID, params.Limit, (params.Page-1)*params.Limit)
	err := r.db.Select(&listItem.Items, query1+mainQuery, listQuery...)

	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			listItem.PlaceInfo = make([]PlaceInfo, 0)
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	err = r.db.Get(&listItem.TotalCount, query2+tempQuery2, listQuery[0:len(listQuery)-2]...)
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			listItem.PlaceInfo = make([]PlaceInfo, 0)
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	err = r.db.Select(&listItem.PlaceInfo, query3, params.PlaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			listItem.PlaceInfo = make([]PlaceInfo, 0)
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

func (r repo) GetListItemAdminWithPagination(params ListItemRequest) (*ListItem, error) {
	var listItem ListItem
	listItem.Items = make([]Item, 0)
	listItem.TotalCount = 0

	query := `
	SELECT i.id, i.name, i.image, i.price, i.description
	FROM items i, places p
	WHERE i.place_id = p.id AND p.user_id = $1 AND i.is_active = TRUE LIMIT $2 OFFSET $3
	`

	err := r.db.Select(&listItem.Items, query, params.UserID, params.Limit, (params.Page-1)*params.Limit)

	if err != nil {
		if err == sql.ErrNoRows {
			listItem.Items = make([]Item, 0)
			listItem.TotalCount = 0
			return &listItem, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = "SELECT COUNT(i.id) FROM items i, places p WHERE i.place_id = p.id AND p.user_id = $1 AND i.is_active = TRUE"

	err = r.db.Get(&listItem.TotalCount, query, params.UserID)
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

func (r repo) DeleteItemAdminByID(itemID int) error {
	query := `
		UPDATE items
		SET is_active = FALSE
		WHERE items.id = $1;
	`

	_, err := r.db.Exec(query, itemID)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

func (r repo) UpdateItem(ID int, item Item) error {
	query := `
		UPDATE items
		SET name=$1, image=$2, description=$3, price=$4, updated_at=now()
		WHERE id=$5
	`

	_, err := r.db.Exec(query, item.Name, item.Image, item.Description, item.Price, ID)
	if err != nil {
		logrus.Errorf("error executing query: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("item not found: %w", ErrNotFound)
		}
		return fmt.Errorf("failed to execute query: %w", ErrInternalServerError)
	}
	return nil
}

func (r repo) CreateItem(userID int, item Item) error {
	query := `
		INSERT INTO items (name, image, description, price, place_id)
		SELECT $1, $2, $3, $4, places.id
		FROM places
		WHERE places.user_id = $5
	`

	_, err := r.db.Exec(query, item.Name, item.Image, item.Description, item.Price, userID)
	if err != nil {
		logrus.Errorf("error executing query: %v", err)
		return fmt.Errorf("failed to execute query: %w", ErrInternalServerError)
	}

	return nil
}
