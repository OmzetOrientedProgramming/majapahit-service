package businessadmin

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service are interface can be used by service
type Service interface {
	GetBalanceDetail(int) (*BalanceDetail, error)
	GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, *util.Pagination, error)
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

	balanceDetail.LatestDisbursementDate = latestDisbursement.Date

	return balanceDetail, nil
}

func (s *service) GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, *util.Pagination, error) {
	var errorList []string
	var listTransaction *ListTransaction
	var err error

	if params.Page == 0 {
		params.Page = util.DefaultPage
	}

	if params.Limit == 0 {
		params.Limit = util.DefaultLimit
	}

	if params.Limit > util.MaxLimit {
		errorList = append(errorList, "limit should be 1 - 100")
	}

	if params.Path == "" {
		errorList = append(errorList, "path is required for pagination")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	listTransaction, err = s.repo.GetListTransactionsHistoryWithPagination(params)

	if err != nil {
		return nil, nil, err
	}

	pagination := util.GeneratePagination(listTransaction.TotalCount, params.Limit, params.Page, params.Path)
	return listTransaction, &pagination, err
}
