package item

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetListItem(t *testing.T) {
	listItemExpected := &ListItem{
		Items: []Item {
			{
				ID:         	1,
				Name:        	"test",
				Price:			10000,
				Description: 	"test",
			},
			{
				ID:          	2,
				Name:        	"test",
				Price:			10000,
				Description: 	"test",
			},
		},
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
		NewRows([]string{"id", "name", "price", "description"}).
		AddRow(listItemExpected.Items[0].ID,
			listItemExpected.Items[0].Name,
			listItemExpected.Items[0].Price,
			listItemExpected.Items[0].Description).
		AddRow(listItemExpected.Items[1].ID,
			listItemExpected.Items[1].Name,
			listItemExpected.Items[1].Price,
			listItemExpected.Items[1].Description)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, price, description FROM items WHERE place_id = $1 ")).
		WithArgs(1).
		WillReturnRows(rows)

	// Test
	listItemResult, err := repoMock.GetListItem(1)
	assert.Equal(t, listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}
