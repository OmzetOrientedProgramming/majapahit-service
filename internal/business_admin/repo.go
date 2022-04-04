package businessadmin

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Repo will contain all the function that can be used by repo
type Repo interface {
	GetPlaceIDByUserID(int) (int, error)
	GetLatestDisbursement(int) (*DisbursementDetail, error)
	GetBalance(int) (*BalanceDetail, error)
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

	query := "SELECT date, amount, status FROM disbursements WHERE (place_id = $1 AND status = 1) ORDER BY date DESC LIMIT 1"
	err := r.db.Get(&result, query, placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
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
