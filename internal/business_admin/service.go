package businessadmin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service are interface can be used by service
type Service interface {
	GetBalanceDetail(int) (*BalanceDetail, error)
	GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, *util.Pagination, error)
	CreateDisbursement(int, float64) (*CreateDisbursementResponse, error)
	DisbursementCallbackFromXendit(params DisbursementCallback) error
	GetTransactionHistoryDetail(int) (*TransactionHistoryDetail, error)
	GetPlaceDetail(userID int) (*PlaceDetail, error)
}

type service struct {
	repo          Repo
	xenditService xendit.Service
	placeService 	place.Service
}

// NewService create new service
func NewService(repo Repo, xenditService xendit.Service, placeService place.Service) Service {
	return &service{
		repo:          repo,
		xenditService: xenditService,
		placeService: placeService,
	}
}

func (s *service) DisbursementCallbackFromXendit(params DisbursementCallback) error {
	userID, err := strconv.Atoi(params.ExternalID)
	if err != nil {
		return errors.Wrap(ErrInputValidationError, "external id is not valid")
	}

	switch params.Status {
	case util.XenditDisbursementCompletedString:
		currentBalance, err := s.repo.GetBalance(userID)
		if err != nil {
			return err
		}

		newBalance := currentBalance.Balance - params.Amount

		err = s.repo.UpdateBalance(newBalance, userID)
		if err != nil {
			return err
		}

		err = s.repo.UpdateDisbursementStatusByXenditID(util.XenditDisbursementCompleted, params.ID)
		if err != nil {
			return err
		}

		return nil
	case util.XenditDisbursementFailedString:
		err = s.repo.UpdateDisbursementStatusByXenditID(util.XenditDisbursementFailed, params.ID)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.Wrap(ErrInputValidationError, "status must be COMPLETED or FAILED")
	}
}

func (s *service) CreateDisbursement(userID int, amount float64) (*CreateDisbursementResponse, error) {
	var errorList []string

	if userID <= 0 {
		errorList = append(errorList, "userID must be positive integer")
	}

	if amount <= 0 {
		errorList = append(errorList, "amount must be positive integer")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	businessAdminInfo, err := s.repo.GetBusinessAdminInformation(userID)
	if err != nil {
		return nil, err
	}

	latestDisbursement, err := s.repo.GetLatestDisbursement(businessAdminInfo.PlaceID)
	if err != nil {
		return nil, err
	}

	currentDateTime := time.Now()
	oneMonthAgo := currentDateTime.AddDate(0, -1, 0)
	if !latestDisbursement.Date.Before(oneMonthAgo) {
		return nil, errors.Wrap(ErrInputValidationError, "disbursement can only be done once a month ")
	}

	amount -= util.XenditDisbursementFee + (util.XenditDisbursementFee * util.XenditVATPercentage)

	xenditDisbursementParams := xendit.CreateDisbursementParams{
		ID:                businessAdminInfo.ID,
		BankAccountName:   businessAdminInfo.BankAccountName,
		BankAccountNumber: businessAdminInfo.BankAccountNumber,
		Amount:            amount,
		Description:       fmt.Sprintf("Disbursement by %s", businessAdminInfo.Name),
		Email:             []string{businessAdminInfo.Email},
	}

	createXenditDisbursement, err := s.xenditService.CreateDisbursement(xenditDisbursementParams)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	disbursement := DisbursementDetail{
		PlaceID:  businessAdminInfo.PlaceID,
		Date:     currentDateTime,
		XenditID: createXenditDisbursement.ID,
		Amount:   createXenditDisbursement.Amount,
		Status:   util.XenditDisbursementPending,
	}

	ID, err := s.repo.SaveDisbursement(disbursement)
	if err != nil {
		return nil, err
	}

	resp := CreateDisbursementResponse{
		ID:        ID,
		CreatedAt: disbursement.Date,
		Amount:    disbursement.Amount,
		XenditID:  createXenditDisbursement.ID,
	}

	return &resp, nil
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

	balanceDetail.LatestDisbursementDate = latestDisbursement.Date.String()

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

func (s *service) GetTransactionHistoryDetail(bookingID int) (*TransactionHistoryDetail, error) {
	var errorList []string

	if bookingID <= 0 {
		errorList = append(errorList, "bookingID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	itemsWrapper, err := s.repo.GetItemsWrapper(bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	customerForTrasactionHistoryDetail, err := s.repo.GetCustomerForTransactionHistoryDetail(bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	transactionHistoryDetail, err := s.repo.GetTransactionHistoryDetail(bookingID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	transactionHistoryDetail.CustomerName = customerForTrasactionHistoryDetail.CustomerName
	transactionHistoryDetail.CustomerImage = customerForTrasactionHistoryDetail.CustomerImage
	transactionHistoryDetail.Items = itemsWrapper.Items

	return transactionHistoryDetail, nil
}

func (s *service) GetPlaceDetail(userID int) (*PlaceDetail, error) {
	errorList := []string{}

	if userID <= 0 {
		errorList = append(errorList, "userID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ";"))
	}

	placeID, err := s.repo.GetPlaceIDByUserID(userID)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServerError, err.Error())
	}

	placeDetail, err := s.placeService.GetDetail(placeID)
	if err != nil {
		return nil, err
	}

	var resPlaceDetail PlaceDetail
	resPlaceDetail.ID = placeDetail.ID
	resPlaceDetail.Name = placeDetail.Name
	resPlaceDetail.Image = placeDetail.Image
	resPlaceDetail.Address = placeDetail.Address
	resPlaceDetail.Description = placeDetail.Description
	resPlaceDetail.OpenHour = placeDetail.OpenHour
	resPlaceDetail.CloseHour = placeDetail.CloseHour
	resPlaceDetail.BookingPrice = placeDetail.BookingPrice
	resPlaceDetail.MinSlot = placeDetail.MinSlot
	resPlaceDetail.MaxSlot = placeDetail.MaxSlot
	resPlaceDetail.Capacity = placeDetail.Capacity
	resPlaceDetail.MinIntervalBooking = placeDetail.MinIntervalBooking
	resPlaceDetail.MaxIntervalBooking = placeDetail.MaxIntervalBooking
	resPlaceDetail.AverageRating = placeDetail.AverageRating

	return &resPlaceDetail, nil
}