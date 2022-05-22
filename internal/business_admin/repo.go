package businessadmin

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetPlaceIDByUserID(int) (int, error)
	GetLatestDisbursement(int) (*DisbursementDetail, error)
	GetBalance(int) (*BalanceDetail, error)
	GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, error)
	GetBusinessAdminInformation(userID int) (*InfoForDisbursement, error)
	SaveDisbursement(disbursement DisbursementDetail) (int, error)
	UpdateBalance(float64, int) error
	UpdateDisbursementStatusByXenditID(int, string) error
	GetTransactionHistoryDetail(int) (*TransactionHistoryDetail, error)
	GetItemsWrapper(int) (*ItemsWrapper, error)
	GetCustomerForTransactionHistoryDetail(int) (*CustomerForTrasactionHistoryDetail, error)
	UpdateProfile(EditProfileRequest) error
}

type repo struct {
	db *sqlx.DB
}

// NewRepo used to initialize repo
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) UpdateBalance(newAmount float64, userID int) error {
	query := "UPDATE business_owners SET balance = $1 " +
		"WHERE user_id= $2"

	_, err := r.db.Exec(query, newAmount, userID)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

func (r *repo) UpdateDisbursementStatusByXenditID(newStatus int, xenditID string) error {
	query := "UPDATE disbursements as d SET status = $1 " +
		"WHERE xendit_id = $2"

	_, err := r.db.Exec(query, newStatus, xenditID)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}

func (r *repo) GetBusinessAdminInformation(userID int) (*InfoForDisbursement, error) {
	var disbursementInfo InfoForDisbursement

	query := "SELECT u.id, u.name, u.email, b.bank_account_number, b.bank_account_name, p.id as place_id " +
		"FROM users as u " +
		"JOIN business_owners as b ON b.user_id = u.id " +
		"JOIN places as p ON p.user_id = u.id " +
		"WHERE u.id = $1"

	err := r.db.Get(&disbursementInfo, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &disbursementInfo, nil
}

func (r *repo) SaveDisbursement(disbursement DisbursementDetail) (int, error) {
	var lastInsertedID int

	query := "INSERT INTO disbursements (place_id, date, xendit_id, amount, status) " +
		"VALUES ($1, $2, $3, $4, $5) " +
		"RETURNING id"

	if err := r.db.QueryRow(query, disbursement.PlaceID, disbursement.Date, disbursement.XenditID, disbursement.Amount, disbursement.Status).Scan(&lastInsertedID); err != nil {
		return 0, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return lastInsertedID, nil
}

func (r *repo) GetPlaceIDByUserID(userID int) (int, error) {
	var result int

	query := "SELECT id FROM places WHERE user_id = $1"
	err := r.db.Get(&result, query, userID)
	if err != nil {
		return 0, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return result, nil
}

func (r *repo) GetLatestDisbursement(placeID int) (*DisbursementDetail, error) {
	var result DisbursementDetail

	query := "SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND (status = 0 OR status = 1)) ORDER BY date DESC LIMIT 1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		if err == sql.ErrNoRows {
			result = DisbursementDetail{
				Date:   time.Time{},
				Amount: 0,
				Status: 1,
			}
		} else {
			return nil, errors.Wrap(ErrInternalServerError, err.Error())
		}
	}

	return &result, nil
}

func (r *repo) GetBalance(userID int) (*BalanceDetail, error) {
	var result BalanceDetail

	query := "SELECT balance FROM business_owners INNER JOIN users ON users.id = business_owners.user_id WHERE business_owners.user_id = $1"
	err := r.db.Get(&result, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &result, nil
}

func (r *repo) GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, error) {
	var listTransaction ListTransaction
	listTransaction.Transactions = make([]Transaction, 0)
	listTransaction.TotalCount = 0

	query := `
	SELECT b.id, u.name, u.image, b.total_price, b.date
	FROM bookings b, users u, places p
	WHERE b.place_id = p.id AND p.user_id = $1 AND b.user_id = u.id AND b.status = 3 
	ORDER BY b.date DESC LIMIT $2 OFFSET $3
	`

	err := r.db.Select(&listTransaction.Transactions, query, params.UserID, params.Limit, (params.Page-1)*params.Limit)

	if err != nil {
		if err == sql.ErrNoRows {
			listTransaction.Transactions = make([]Transaction, 0)
			listTransaction.TotalCount = 0
			return &listTransaction, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	query = "SELECT COUNT(b.id) FROM bookings b, places p WHERE b.place_id = p.id AND p.user_id = $1 AND b.status = 3"

	err = r.db.Get(&listTransaction.TotalCount, query, params.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			listTransaction.Transactions = make([]Transaction, 0)
			listTransaction.TotalCount = 0
			return &listTransaction, nil
		}
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &listTransaction, nil
}

func (r *repo) GetTransactionHistoryDetail(bookingID int) (*TransactionHistoryDetail, error) {
	var transactionHistoryDetail TransactionHistoryDetail

	query := `
		SELECT date, start_time, end_time, total_price, capacity
		FROM bookings 
		WHERE id = $1
	`

	err := r.db.Get(&transactionHistoryDetail, query, bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &transactionHistoryDetail, nil
}

func (r *repo) GetItemsWrapper(bookingID int) (*ItemsWrapper, error) {
	var itemsWrapper ItemsWrapper
	itemsWrapper.Items = make([]ItemDetail, 0)

	query := `
		SELECT i.name, bi.qty, i.price
		FROM bookings b
		INNER JOIN booking_items bi
		ON b.id = bi.booking_id
		INNER JOIN items i
		ON bi.item_id = i.id
		WHERE b.id = $1
	`
	err := r.db.Select(&itemsWrapper.Items, query, bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &itemsWrapper, nil
}

func (r *repo) GetCustomerForTransactionHistoryDetail(bookingID int) (*CustomerForTrasactionHistoryDetail, error) {
	var customer CustomerForTrasactionHistoryDetail

	query := `
		SELECT u.name, u.image
		FROM bookings b
		INNER JOIN users u
		ON b.user_id = u.id
		WHERE b.id = $1
	`
	err := r.db.Get(&customer, query, bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	return &customer, nil
}

func (r *repo) UpdateProfile(businessAdminProfile EditProfileRequest) error {
	query := `
		UPDATE places
		SET name = $1,
			description = $2
		WHERE user_id = $3
  	`

	_, err := r.db.Exec(query, businessAdminProfile.Name, businessAdminProfile.Description, businessAdminProfile.UserID)
	if err != nil {
		return errors.Wrap(ErrInternalServerError, err.Error())
	}

	return nil
}
