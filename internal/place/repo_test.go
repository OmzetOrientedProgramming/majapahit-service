package place

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetPlaceDetailSuccess(t *testing.T) {
	placeId := 1
	placeDetailExpected := &PlaceDetail{
		ID:          1,
		Name:        "test_name_place",
		Image:       "test_image_place",
		Distance:    200,
		Address:     "test_address_place",
		Description: "test_description_place",
		OpenHour:    "08:00",
		CloseHour:   "16:00",
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
		NewRows([]string{"id", "name", "image", "distance", "address", "description", "open_hour", "close_hour"}).
		AddRow(
			placeDetailExpected.ID,
			placeDetailExpected.Name,
			placeDetailExpected.Image,
			placeDetailExpected.Distance,
			placeDetailExpected.Address,
			placeDetailExpected.Description,
			placeDetailExpected.OpenHour,
			placeDetailExpected.CloseHour,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, distance, address, description, open_hour, close_hour FROM places WHERE id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	// Test
	placeDetailRetrieve, err := repoMock.GetPlaceDetail(placeId)
	assert.Equal(t, placeDetailExpected, placeDetailRetrieve)
	assert.NotNil(t, placeDetailRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetPlaceDetailInternalServerError(t *testing.T) {
	placeId := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, distance, address, description, open_hour, close_hour, rating FROM places WHERE id = $1")).
		WithArgs(placeId).
		WillReturnError(sql.ErrTxDone)

	// Test
	placeDetailRetrieve, err := repoMock.GetPlaceDetail(placeId)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailRetrieve)
}
