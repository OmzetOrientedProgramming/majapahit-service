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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND (status = 0 OR status = 1)) ORDER BY date DESC LIMIT 1")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND (status = 0 OR status = 1)) ORDER BY date DESC LIMIT 1")).
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND (status = 0 OR status = 1)) ORDER BY date DESC LIMIT 1")).
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

func TestRepo_GetListTransactionsHistoryWithPaginationSuccess(t *testing.T) {
	listTransactionExpected := &ListTransaction{
		Transactions: []Transaction{
			{
				ID:    1,
				Name:  "test name",
				Image: "test image",
				Price: 10000,
				Date:  "test date",
			},
			{
				ID:    2,
				Name:  "test name",
				Image: "test image",
				Price: 10000,
				Date:  "test date",
			},
		},
		TotalCount: 10,
	}

	params := ListTransactionRequest{
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
		NewRows([]string{"id", "name", "image", "total_price", "date"}).
		AddRow(listTransactionExpected.Transactions[0].ID,
			listTransactionExpected.Transactions[0].Name,
			listTransactionExpected.Transactions[0].Image,
			listTransactionExpected.Transactions[0].Price,
			listTransactionExpected.Transactions[0].Date).
		AddRow(listTransactionExpected.Transactions[1].ID,
			listTransactionExpected.Transactions[1].Name,
			listTransactionExpected.Transactions[1].Image,
			listTransactionExpected.Transactions[1].Price,
			listTransactionExpected.Transactions[1].Date)
	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT b.id, u.name, u.image, b.total_price, b.date
	FROM bookings b, users u, places p
	WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)

	rows = mock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnRows(rows)

	// Test
	listTransactionResult, err := repoMock.GetListTransactionsHistoryWithPagination(params)
	assert.Equal(t, listTransactionExpected, listTransactionResult)
	assert.NotNil(t, listTransactionResult)
	assert.NoError(t, err)
}

func TestRepo_GetListTransactionsHistoryWithPaginationError(t *testing.T) {
	params := ListTransactionRequest{
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
		SELECT b.id, u.name, u.image, b.total_price, b.date
		FROM bookings b, users u, places p
		WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrTxDone)

	// Test
	listTransactionsResult, err := repoMock.GetListTransactionsHistoryWithPagination(params)
	assert.Nil(t, listTransactionsResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListTransactionsHistoryWithPaginationCountError(t *testing.T) {
	params := ListTransactionRequest{
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
		NewRows([]string{"id", "name", "image", "total_price", "date"}).
		AddRow("1", "test name", "image", 10, "date")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT b.id, u.name, u.image, b.total_price, b.date
		FROM bookings b, users u, places p
		WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnError(sql.ErrConnDone)

	// Test
	listTransactionsResult, err := repoMock.GetListTransactionsHistoryWithPagination(params)
	assert.Nil(t, listTransactionsResult)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))

}

func TestRepo_GetListTransactionsHistoryWithPaginationEmpty(t *testing.T) {
	listTransactionsExpected := &ListTransaction{
		Transactions: make([]Transaction, 0),
	}

	params := ListTransactionRequest{
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
		SELECT b.id, u.name, u.image, b.total_price, b.date
		FROM bookings b, users u, places p
		WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnError(sql.ErrNoRows)

	// Test
	listTransactionsResult, err := repoMock.GetListTransactionsHistoryWithPagination(params)
	assert.Equal(t, listTransactionsExpected, listTransactionsResult)
	assert.NotNil(t, listTransactionsResult)
	assert.NoError(t, err)

}

func TestRepo_GetListTransactionsHistoryWithPaginationCountEmpty(t *testing.T) {
	listTransactionsExpected := &ListTransaction{
		Transactions: make([]Transaction, 0),
		TotalCount:   0,
	}

	params := ListTransactionRequest{
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
		NewRows([]string{"id", "name", "image", "total_price", "date"}).
		AddRow("1", "name", "image", 10, "date")
	mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT b.id, u.name, u.image, b.total_price, b.date
	FROM bookings b, users u, places p
	WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id LIMIT $2 OFFSET $3`)).
		WithArgs(params.UserID, params.Limit, (params.Page-1)*params.Limit).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1")).
		WithArgs(params.UserID).
		WillReturnError(sql.ErrNoRows)

	// Test
	listTransactionsResult, err := repoMock.GetListTransactionsHistoryWithPagination(params)
	assert.Equal(t, listTransactionsExpected, listTransactionsResult)
	assert.NotNil(t, listTransactionsResult)
	assert.NoError(t, err)
}

func TestRepo_GetTransactionHistoryDetailSuccess(t *testing.T) {
	bookingID := 1
	transactionHistoryDetailExpected := &TransactionHistoryDetail{
		Date:           "27 Oktober 2021",
		StartTime:      "20-00",
		EndTime:        "21-00",
		TotalPriceItem: 500000,
		Capacity:       5,
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
		NewRows([]string{"date", "start_time", "end_time", "total_price", "capacity"}).
		AddRow(
			transactionHistoryDetailExpected.Date,
			transactionHistoryDetailExpected.StartTime,
			transactionHistoryDetailExpected.EndTime,
			transactionHistoryDetailExpected.TotalPriceItem,
			transactionHistoryDetailExpected.Capacity,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT date, start_time, end_time, total_price, capacity 
			FROM bookings 
			WHERE id = $1
		`)).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	transactionHistoryDetailRetrieved, err := repoMock.GetTransactionHistoryDetail(bookingID)
	assert.Equal(t, transactionHistoryDetailExpected, transactionHistoryDetailRetrieved)
	assert.NotNil(t, transactionHistoryDetailRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetTransactionHistoryDetailInternalServerError(t *testing.T) {
	bookingID := 1

	// mockDB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT date, start_time, end_time, total_price capacity 
									   FROM bookings 
									   WHERE id = $1`)).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	transactionHistoryDetailRetrieved, err := repoMock.GetTransactionHistoryDetail(bookingID)
	assert.Nil(t, transactionHistoryDetailRetrieved)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetItemsWrapperSuccess(t *testing.T) {
	bookingID := 1
	itemsWrapperExpected := &ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "Sample Name 1",
				Qty:   10,
				Price: 75000.0,
			},
			{
				Name:  "Sample Name 2",
				Qty:   5,
				Price: 20000.0,
			},
		},
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
		NewRows([]string{"name", "qty", "price"}).
		AddRow(
			itemsWrapperExpected.Items[0].Name,
			itemsWrapperExpected.Items[0].Qty,
			itemsWrapperExpected.Items[0].Price,
		).
		AddRow(
			itemsWrapperExpected.Items[1].Name,
			itemsWrapperExpected.Items[1].Qty,
			itemsWrapperExpected.Items[1].Price,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT i.name, bi.qty, i.price
			FROM bookings b
			INNER JOIN booking_items bi
			ON b.id = bi.booking_id
			INNER JOIN items i
			ON bi.item_id = i.id
			WHERE b.id = $1
		`)).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	itemsWrapperRetrieved, err := repoMock.GetItemsWrapper(bookingID)
	assert.Equal(t, itemsWrapperExpected, itemsWrapperRetrieved)
	assert.NotNil(t, itemsWrapperRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetItemsWrapperInternalServerError(t *testing.T) {
	bookingID := 1

	// mockDB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)
	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT i.name, bi.qty, i.price
			FROM bookings b
			INNER JOIN booking_items bi
			ON b.id = bi.booking_id
			INNER JOIN items i
			ON bi.item_id = i.id
			WHERE b.id = $1
		`)).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	itemsWrapperRetrieved, err := repoMock.GetItemsWrapper(bookingID)
	assert.Nil(t, itemsWrapperRetrieved)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}

func TestRepo_GetCustomerForTransactionHistoryDetailSuccess(t *testing.T) {
	bookingID := 1
	customerExpected := &CustomerForTrasactionHistoryDetail{
		CustomerName:  "customer_name_1",
		CustomerImage: "customer_image_1",
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
		NewRows([]string{"name", "image"}).
		AddRow(
			customerExpected.CustomerName,
			customerExpected.CustomerImage,
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT u.name, u.image
			FROM bookings b
			INNER JOIN users u
			ON b.user_id = u.id
			WHERE b.id = $1
		`)).
		WithArgs(bookingID).
		WillReturnRows(rows)

	// Test
	customerRetrieved, err := repoMock.GetCustomerForTransactionHistoryDetail(bookingID)
	assert.Equal(t, customerExpected, customerRetrieved)
	assert.NotNil(t, customerRetrieved)
	assert.NoError(t, err)
}

func TestRepo_GetCustomerForTransactionHistoryDetailInternalServerError(t *testing.T) {
	bookingID := 1

	// mockDB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Expectation
	repoMock := NewRepo(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT u.name, u.image
			FROM bookings b
			INNER JOIN users u
			ON b.user_id = u.id
			WHERE b.id = $1
		`)).
		WithArgs(bookingID).
		WillReturnError(sql.ErrTxDone)

	// Test
	customerRetrieved, err := repoMock.GetCustomerForTransactionHistoryDetail(bookingID)
	assert.Nil(t, customerRetrieved)
	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
}
