package place

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetPlaceDetailSuccess(t *testing.T) {
	placeId := 1
	placeDetailExpected := &PlaceDetail{
		ID:            1,
		Name:          "test_name_place",
		Image:         "test_image_place",
		Distance:      200,
		Address:       "test_address_place",
		Description:   "test_description_place",
		OpenHour:      "08:00",
		CloseHour:     "16:00",
		AverageRating: 3.5,
		ReviewCount:   15,
		Reviews: []UserReview{
			{
				User:    "test_user_1",
				Rating:  5,
				Content: "test_content_user_1",
			},
			{
				User:    "test_user_2",
				Rating:  4,
				Content: "test_content_user_2",
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
		NewRows([]string{"id", "name", "image", "distance", "address", "description", "open_hour", "close_hour", "rating"}).
		AddRow(
			placeDetailExpected.ID,
			placeDetailExpected.Name,
			placeDetailExpected.Image,
			placeDetailExpected.Distance,
			placeDetailExpected.Address,
			placeDetailExpected.Description,
			placeDetailExpected.OpenHour,
			placeDetailExpected.CloseHour,
			placeDetailExpected.AverageRating,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, distance, address, description, open_hour, close_hour, rating FROM places WHERE id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"name", "rating", "content"}).
		AddRow(
			placeDetailExpected.Reviews[0].User,
			placeDetailExpected.Reviews[0].Rating,
			placeDetailExpected.Reviews[0].Content).
		AddRow(
			placeDetailExpected.Reviews[1].User,
			placeDetailExpected.Reviews[1].Rating,
			placeDetailExpected.Reviews[1].Content,
		)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(15)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT users.name, reviews.rating, reviews.content FROM reviews LEFT JOIN users WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	// Test
	placeDetailRetrieve, err := repoMock.GetPlaceDetail(placeId)
	assert.Equal(t, placeDetailExpected, placeDetailRetrieve)
	assert.NotNil(t, placeDetailRetrieve)
	assert.NoError(t, err)

}
