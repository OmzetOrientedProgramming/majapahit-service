package businessadmin

import (
	"strings"

	"github.com/pkg/errors"
)

// Service are interface can be used by service
type Service interface {
	GetBalanceDetail(int) (*BalanceDetail, error)
}

type service struct {
	repo Repo
}

// NewService create new service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetBalanceDetail(userID int) (*BalanceDetail, error) {
	var errorList []string

	if userID <= 0 {
		errorList = append(errorList, "userID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	placeID, err := s.repo.GetPlaceIDByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	latestDisbursement, err := s.repo.GetLatestDisbursement(placeID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	balanceDetail, err := s.repo.GetBalance(userID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	if latestDisbursement.Status == 1 {
		balanceDetail.Balance = balanceDetail.Balance - latestDisbursement.Amount
	}

	balanceDetail.LatestDisbursementDate = latestDisbursement.Date

	return balanceDetail, nil
}
