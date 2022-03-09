package place

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

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
