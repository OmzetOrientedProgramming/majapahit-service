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

func TestRepo_GetUserReviewForPlaceDetailSuccess(t *testing.T) {
	placeId := 1
	expectedAverageRatingAndReviews := &AverageRatingAndReviews{
		AverageRating: 3.50,
		ReviewCount:   30,
		Reviews: []UserReview{
			{
				User:    "test_user_1",
				Rating:  4.50,
				Content: "test_review_content_1",
			},
			{
				User:    "test_user_2",
				Rating:  5,
				Content: "test_review_content_2",
			},
		},
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

	rows := mock.NewRows([]string{"count_review"}).AddRow(30)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"user", "rating", "content"}).
		AddRow(
			expectedAverageRatingAndReviews.Reviews[0].User,
			expectedAverageRatingAndReviews.Reviews[0].Rating,
			expectedAverageRatingAndReviews.Reviews[0].Content).
		AddRow(
			expectedAverageRatingAndReviews.Reviews[1].User,
			expectedAverageRatingAndReviews.Reviews[1].Rating,
			expectedAverageRatingAndReviews.Reviews[1].Content)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2")).
		WithArgs(placeId).
		WillReturnRows(rows)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeId)
	assert.Equal(t, expectedAverageRatingAndReviews, retrivedAverageRatingAndReviews)
	assert.NotNil(t, retrivedAverageRatingAndReviews)
	assert.NoError(t, err)
}

func TestRepo_GetUserReviewForPlaceDetailCountReviewInternalServerError(t *testing.T) {
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeId)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}

func TestRepo_GetUserReviewForPlaceDetailSumRatingInternalServerError(t *testing.T) {
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

	rows := mock.NewRows([]string{"count_review"}).AddRow(30)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeId)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}

func TestRepo_GetUserReviewForPlaceDetailInternalServerError(t *testing.T) {
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

	rows := mock.NewRows([]string{"count_review"}).AddRow(30)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) as count_review FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeId).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2")).
		WithArgs(placeId).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeId)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}
