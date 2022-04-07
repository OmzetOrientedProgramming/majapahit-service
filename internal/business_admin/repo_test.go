package businessadmin

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRepo_GetLatestDisbursementSuccess(t *testing.T) {
	placeID := 1
	latestDateExpected := &DisbursementDetail{
		Date:   "27 Oktober 2021",
		Amount: 150000,
		Status: 0,
	}

	// mockDB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	rows := mock.
		NewRows([]string{"date", "amount", "status"}).
		AddRow(
			latestDateExpected.Date,
			latestDateExpected.Amount,
			latestDateExpected.Status,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND status = 1) ORDER BY date DESC LIMIT 1")).
		WithArgs(placeID).
		WillReturnRows(rows)

	// Test
	latestDateRetrieved, err := repoMock.GetLatestDisbursement(placeID)
	assert.Equal(t, latestDateExpected, latestDateRetrieved)
	assert.NotNil(t, latestDateRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetLatestDisbursementErrorNoRows(t *testing.T) {
	placeID := 1
	latestDateExpected := &DisbursementDetail{
		Date:   "-",
		Amount: 0,
		Status: 1,
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND status = 1) ORDER BY date DESC LIMIT 1")).
		WithArgs(placeID).
		WillReturnError(sql.ErrNoRows)

	// Test
	latestDateRetrieved, err := repoMock.GetLatestDisbursement(placeID)
	assert.Equal(t, latestDateExpected, latestDateRetrieved)
	assert.NotNil(t, latestDateRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetLatestDisbursementInternalServerError(t *testing.T) {
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND status = 1) ORDER BY date DESC LIMIT 1")).
		WithArgs(placeID).
		WillReturnError(sql.ErrTxDone)

	// Test
	latestDateRetrieved, err := repoMock.GetLatestDisbursement(placeID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, latestDateRetrieved)
}

func TestRepo_GetPlaceIDByUserIDSuccess(t *testing.T) {
	userID := 1
	placeIDExpected := 3

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
		NewRows([]string{"id"}).
		AddRow(
			placeIDExpected,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM places WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnRows(rows)

	// Test
	placeIDRetrieved, err := repoMock.GetPlaceIDByUserID(userID)
	assert.Equal(t, placeIDExpected, placeIDRetrieved)
	assert.NotNil(t, placeIDRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetPlaceIDByUserIDInternalServerError(t *testing.T) {
	userID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM places WHERE user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

	// Test
	placeIDRetrieved, err := repoMock.GetPlaceIDByUserID(userID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Equal(t, 0, placeIDRetrieved)
}

func TestRepo_GetBalanceSuccess(t *testing.T) {
	userID := 1
	balanceDetailExpected := &BalanceDetail{
		Balance: 10000000,
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
		NewRows([]string{"balance"}).AddRow(balanceDetailExpected.Balance)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM business_owners INNER JOIN users ON users.id = business_owners.user_id WHERE business_owners.user_id = $1")).
		WithArgs(userID).
		WillReturnRows(rows)

	// Test
	balanceDetailRetrieved, err := repoMock.GetBalance(userID)
	assert.Equal(t, balanceDetailExpected, balanceDetailRetrieved)
	assert.NotNil(t, balanceDetailRetrieved)
	assert.NoError(t, err)
}
func TestRepo_GetBalanceInternalServerError(t *testing.T) {
	userID := 1

	// Mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT balance FROM business_owners INNER JOIN users ON users.id = business_owners.user_id WHERE business_owners.user_id = $1")).
		WithArgs(userID).
		WillReturnError(sql.ErrTxDone)

	// Test
	balanceDetailRetrieved, err := repoMock.GetBalance(userID)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, balanceDetailRetrieved)
}
