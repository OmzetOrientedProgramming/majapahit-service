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

func TestRepo_GetDetailSuccess(t *testing.T) {
	placeID := 1
	placeDetailExpected := &Detail{
		ID:           1,
		Name:         "test_name_place",
		Image:        "test_image_place",
		Distance:     200,
		Address:      "test_address_place",
		Description:  "test_description_place",
		OpenHour:     "08:00",
		CloseHour:    "16:00",
		BookingPrice: 15000,
		MinSlot:      2,
		MaxSlot:      5,
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
		NewRows([]string{"id", "name", "image", "distance", "address", "description", "open_hour", "close_hour", "booking_price", "min_slot_booking", "max_slot_booking"}).
		AddRow(
			placeDetailExpected.ID,
			placeDetailExpected.Name,
			placeDetailExpected.Image,
			placeDetailExpected.Distance,
			placeDetailExpected.Address,
			placeDetailExpected.Description,
			placeDetailExpected.OpenHour,
			placeDetailExpected.CloseHour,
			placeDetailExpected.BookingPrice,
			placeDetailExpected.MinSlot,
			placeDetailExpected.MaxSlot,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, distance, address, description, open_hour, close_hour, booking_price, min_slot_booking, max_slot_booking FROM places WHERE id = $1")).
		WithArgs(placeID).
		WillReturnRows(rows)

	// Test
	placeDetailRetrieve, err := repoMock.GetDetail(placeID)
	assert.Equal(t, placeDetailExpected, placeDetailRetrieve)
	assert.NotNil(t, placeDetailRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetDetailInternalServerError(t *testing.T) {
	placeID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, image, distance, address, description, open_hour, close_hour, booking_price, min_slot_booking, max_slot_booking FROM places WHERE id = $1")).
		WithArgs(placeID).
		WillReturnError(sql.ErrTxDone)

	// Test
	placeDetailRetrieve, err := repoMock.GetDetail(placeID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailRetrieve)
}

func TestRepo_GetUserReviewForDetailSuccess(t *testing.T) {
	placeID := 1
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
		WithArgs(placeID).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeID).
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
		WithArgs(placeID).
		WillReturnRows(rows)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeID)
	assert.Equal(t, expectedAverageRatingAndReviews, retrivedAverageRatingAndReviews)
	assert.NotNil(t, retrivedAverageRatingAndReviews)
	assert.NoError(t, err)
}

func TestRepo_GetUserReviewForDetailCountReviewInternalServerError(t *testing.T) {
	placeID := 1

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
		WithArgs(placeID).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}

func TestRepo_GetUserReviewForDetailSumRatingInternalServerError(t *testing.T) {
	placeID := 1

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
		WithArgs(placeID).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeID).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}

func TestRepo_GetUserReviewForDetailInternalServerError(t *testing.T) {
	placeID := 1

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
		WithArgs(placeID).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"sum_rating"}).AddRow(105)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(rating) as sum_rating FROM reviews WHERE place_id = $1")).
		WithArgs(placeID).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT users.name as user, reviews.rating as rating, reviews.content as content FROM reviews LEFT JOIN users ON reviews.user_id = users.id WHERE reviews.place_id = $1 LIMIT 2")).
		WithArgs(placeID).
		WillReturnError(sql.ErrTxDone)

	// Test
	retrivedAverageRatingAndReviews, err := repoMock.GetAverageRatingAndReviews(placeID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, retrivedAverageRatingAndReviews)
}

func TestRepo_GetPlacesListWithPaginationSuccess(t *testing.T) {
	placeListExpected := &PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test",
				Description: "test",
				Image:       "test/image.png",
			},
			{
				ID:          2,
				Name:        "test 2",
				Description: "test 2",
				Image:       "test/image.png",
			},
		},
		TotalCount: 10,
	}

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
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
		NewRows([]string{"id", "name", "description", "address", "image"}).
		AddRow(placeListExpected.Places[0].ID,
			placeListExpected.Places[0].Name,
			placeListExpected.Places[0].Description,
			placeListExpected.Places[0].Address,
			placeListExpected.Places[0].Image).
		AddRow(placeListExpected.Places[1].ID,
			placeListExpected.Places[1].Name,
			placeListExpected.Places[1].Description,
			placeListExpected.Places[1].Address,
			placeListExpected.Places[1].Image)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2")).
		WithArgs(params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM places")).
		WillReturnRows(rows)

	// Test
	placeListRetrieve, err := repoMock.GetPlacesListWithPagination(params)
	assert.Equal(t, placeListExpected, placeListRetrieve)
	assert.NotNil(t, placeListRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetPlacesListWithPaginationEmpty(t *testing.T) {
	placeListExpected := &PlacesList{
		Places:     make([]Place, 0),
		TotalCount: 0,
	}

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2")).
		WithArgs(params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	placeListRetrieve, err := repoMock.GetPlacesListWithPagination(params)
	assert.Equal(t, placeListExpected, placeListRetrieve)
	assert.NotNil(t, placeListRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetPlacesListWithPaginationEmptyWhenCount(t *testing.T) {
	placeListExpected := &PlacesList{
		Places:     make([]Place, 0),
		TotalCount: 0,
	}

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
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
		NewRows([]string{"id", "name", "description", "address"}).
		AddRow("1", "test name", "description", "address")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2")).
		WithArgs(params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM places")).
		WillReturnError(sql.ErrNoRows)

	// Test
	placeListRetrieve, err := repoMock.GetPlacesListWithPagination(params)
	assert.Equal(t, placeListExpected, placeListRetrieve)
	assert.NotNil(t, placeListRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetPlacesListWithPaginationError(t *testing.T) {
	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
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
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2")).
		WithArgs(params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	placeListRetrieve, err := repoMock.GetPlacesListWithPagination(params)
	assert.Nil(t, placeListRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetPlacesListWithPaginationErrorWhenCount(t *testing.T) {
	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
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
		NewRows([]string{"id", "name", "description", "address"}).
		AddRow("1", "test name", "description", "address")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, address, image FROM places LIMIT $1 OFFSET $2")).
		WithArgs(params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM places")).
		WillReturnError(sql.ErrConnDone)

	// Test
	placeListRetrieve, err := repoMock.GetPlacesListWithPagination(params)
	assert.Nil(t, placeListRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetPlaceReviewSuccess(t *testing.T) {
	ratingData := []int{1, 2}
	sumRating := 0
	for _, rating := range ratingData {
		sumRating += rating
	}
	averageRating := float64(sumRating) / float64(len(ratingData))

	expectedResult := PlacesRatingAndReviewCount{
		Rating:      averageRating,
		ReviewCount: len(ratingData),
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
		NewRows([]string{"review_count", "rating"}).
		AddRow(len(ratingData), averageRating)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(rating) as review_count, COALESCE (AVG(rating), 0.0) as rating FROM reviews WHERE place_id = $1")).
		WithArgs(1).
		WillReturnRows(rows)

	// Test
	ratingAndReviewCountRetrieve, err := repoMock.GetPlaceRatingAndReviewCountByPlaceID(1)
	assert.Equal(t, &expectedResult, ratingAndReviewCountRetrieve)
	assert.NotNil(t, ratingAndReviewCountRetrieve)
	assert.NoError(t, err)
}

func TestRepo_GetPlaceReviewFailed(t *testing.T) {
	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Expectation
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(rating) as review_count, COALESCE (AVG(rating), 0.0) as rating FROM reviews WHERE place_id = $1")).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	// Test
	placeListRetrieve, err := repoMock.GetPlaceRatingAndReviewCountByPlaceID(1)
	assert.Nil(t, placeListRetrieve)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}
