package item

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetListItemWwithPaginationSuccess(t *testing.T) {
	listItemExpected := &ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
		},
		TotalCount: 10,
		PlaceInfo: []PlaceInfo{
			{
				Name:  "test",
				Image: "test",
			},
		},
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow(listItemExpected.Items[0].ID,
			listItemExpected.Items[0].Name,
			listItemExpected.Items[0].Image,
			listItemExpected.Items[0].Price,
			listItemExpected.Items[0].Description).
		AddRow(listItemExpected.Items[1].ID,
			listItemExpected.Items[1].Name,
			listItemExpected.Items[1].Image,
			listItemExpected.Items[1].Price,
			listItemExpected.Items[1].Description)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"name", "image"}).
		AddRow(listItemExpected.PlaceInfo[0].Name,
			listItemExpected.PlaceInfo[0].Image)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, image FROM places WHERE id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestRepo_GetListItemWithPaginationError(t *testing.T) {
	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Nil(t, listItemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListItemWithPaginationCountError(t *testing.T) {
	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "test name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrConnDone)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Nil(t, listItemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListItemWithPaginationPlaceInfoError(t *testing.T) {
	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "test name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, image FROM places WHERE id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrConnDone)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Nil(t, listItemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListItemWithPaginationEmpty(t *testing.T) {
	listItemExpected := &ListItem{
		Items:     make([]Item, 0),
		PlaceInfo: make([]PlaceInfo, 0),
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image,  price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)

}

func TestRepo_GetListItemWithPaginationCountEmpty(t *testing.T) {
	listItemExpected := &ListItem{
		Items:      make([]Item, 0),
		TotalCount: 0,
		PlaceInfo:  make([]PlaceInfo, 0),
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrNoRows)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestRepo_GetListItemWithPaginationPlaceInfoEmpty(t *testing.T) {
	listItemExpected := &ListItem{
		Items:      make([]Item, 0),
		TotalCount: 0,
		PlaceInfo:  make([]PlaceInfo, 0),
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 LIMIT $2 OFFSET $3")).
		WithArgs(params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE place_id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, image FROM places WHERE id = $1")).
		WithArgs(params.PlaceID).
		WillReturnError(sql.ErrNoRows)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestRepo_GetListItemWithPaginationByName(t *testing.T) {
	listItemExpected := &ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
		},
		TotalCount: 10,
		PlaceInfo: []PlaceInfo{
			{
				Name:  "test",
				Image: "test",
			},
		},
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
		Name:    "test",
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow(listItemExpected.Items[0].ID,
			listItemExpected.Items[0].Name,
			listItemExpected.Items[0].Image,
			listItemExpected.Items[0].Price,
			listItemExpected.Items[0].Description).
		AddRow(listItemExpected.Items[1].ID,
			listItemExpected.Items[1].Name,
			listItemExpected.Items[1].Image,
			listItemExpected.Items[1].Price,
			listItemExpected.Items[1].Description)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE LOWER(name) LIKE LOWER($1) AND place_id = $2 LIMIT $3 OFFSET $4")).
		WithArgs("%"+params.Name+"%", params.PlaceID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM items WHERE LOWER(name) LIKE LOWER($1) AND place_id = $2")).
		WithArgs("%"+params.Name+"%", params.PlaceID).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"name", "image"}).
		AddRow(listItemExpected.PlaceInfo[0].Name,
			listItemExpected.PlaceInfo[0].Image)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT name, image FROM places WHERE id = $1")).
		WithArgs(params.PlaceID).
		WillReturnRows(rows)

	// Test
	listItemResult, err := repoMock.GetListItemWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestRepo_GetItemByID(t *testing.T) {
	itemExpected := &Item{
		ID:          1,
		Name:        "test",
		Image:       "test",
		Price:       10000,
		Description: "test",
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow(itemExpected.ID,
			itemExpected.Name,
			itemExpected.Image,
			itemExpected.Price,
			itemExpected.Description)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 AND id = $2")).
		WithArgs(10, 1).
		WillReturnRows(rows)

	// Test
	itemResult, err := repoMock.GetItemByID(10, 1)
	assert.Equal(t, itemExpected, itemResult)
	assert.NotNil(t, itemResult)
	assert.NoError(t, err)
}

func TestRepo_GetItemByIDError(t *testing.T) {
	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 AND id = $2")).
		WithArgs(10, 1).
		WillReturnError(sql.ErrTxDone)

	// Test
	itemResult, err := repoMock.GetItemByID(10, 1)
	assert.Nil(t, itemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetItemByIDEmpty(t *testing.T) {
	itemExpected := &Item{}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, price, description FROM items WHERE place_id = $1 AND id = $2")).
		WithArgs(10, 1).
		WillReturnError(sql.ErrNoRows)

	// Test
	itemResult, err := repoMock.GetItemByID(10, 1)
	assert.Equal(t, itemExpected, itemResult)
	assert.NotNil(t, itemResult)
	assert.NoError(t, err)
}

func TestRepo_GetListItemAdminWwithPaginationSuccess(t *testing.T) {
	listItemExpected := &ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Price:       10000,
				Description: "test",
			},
		},
		TotalCount: 10,
	}

	params := ListItemRequest{
		Limit:  10,
		Page:   1,
		UserID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow(listItemExpected.Items[0].ID,
			listItemExpected.Items[0].Name,
			listItemExpected.Items[0].Image,
			listItemExpected.Items[0].Price,
			listItemExpected.Items[0].Description).
		AddRow(listItemExpected.Items[1].ID,
			listItemExpected.Items[1].Name,
			listItemExpected.Items[1].Image,
			listItemExpected.Items[1].Price,
			listItemExpected.Items[1].Description)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT i.id, i.name, i.image, i.price, i.description
		FROM items i, places p
		WHERE i.place_id = p.id AND p.user_id = $1 LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(i.id) FROM items i, places p WHERE i.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnRows(rows)

	// Test
	listItemResult, err := repoMock.GetListItemAdminWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestRepo_GetListItemAdminWithPaginationError(t *testing.T) {
	params := ListItemRequest{
		Limit:  10,
		Page:   1,
		UserID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT i.id, i.name, i.image, i.price, i.description
		FROM items i, places p
		WHERE i.place_id = p.id AND p.user_id = $1 LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	listItemResult, err := repoMock.GetListItemAdminWithPagination(params)
	assert.Nil(t, listItemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListItemAdminWithPaginationCountError(t *testing.T) {
	params := ListItemRequest{
		Limit:  10,
		Page:   1,
		UserID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "test name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT i.id, i.name, i.image, i.price, i.description
		FROM items i, places p
		WHERE i.place_id = p.id AND p.user_id = $1 LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(i.id) FROM items i, places p WHERE i.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnError(sql.ErrConnDone)

	// Test
	listItemResult, err := repoMock.GetListItemAdminWithPagination(params)
	assert.Nil(t, listItemResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListItemAdminWithPaginationEmpty(t *testing.T) {
	listItemExpected := &ListItem{
		Items: make([]Item, 0),
	}

	params := ListItemRequest{
		Limit:  10,
		Page:   1,
		UserID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT i.id, i.name, i.image, i.price, i.description
		FROM items i, places p
		WHERE i.place_id = p.id AND p.user_id = $1 LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	listItemResult, err := repoMock.GetListItemAdminWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)

}

func TestRepo_GetListItemAdminWithPaginationCountEmpty(t *testing.T) {
	listItemExpected := &ListItem{
		Items:      make([]Item, 0),
		TotalCount: 0,
	}

	params := ListItemRequest{
		Limit:  10,
		Page:   1,
		UserID: 1,
	}

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"id", "name", "image", "price", "description"}).
		AddRow("1", "name", "image", 10, "description")
	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT i.id, i.name, i.image, i.price, i.description
	FROM items i, places p
	WHERE i.place_id = p.id AND p.user_id = $1 LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(i.id) FROM items i, places p WHERE i.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnError(sql.ErrNoRows)

	// Test
	listItemResult, err := repoMock.GetListItemAdminWithPagination(params)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}
