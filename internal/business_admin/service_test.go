package businessadmin

import (
	"fmt"
	"testing"
	"time"

	"github.com/tkuchiki/faketime"
	xendit2 "github.com/xendit/xendit-go"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetPlaceIDByUserID(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetLatestDisbursement(placeID int) (*DisbursementDetail, error) {
	args := m.Called(placeID)
	ret := args.Get(0).(DisbursementDetail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetBalance(userID int) (*BalanceDetail, error) {
	args := m.Called(userID)
	ret := args.Get(0).(BalanceDetail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, error) {
	args := m.Called(params)
	ret := args.Get(0).(ListTransaction)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetTransactionHistoryDetail(bookingID int) (*TransactionHistoryDetail, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(TransactionHistoryDetail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetItemsWrapper(bookingID int) (*ItemsWrapper, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(ItemsWrapper)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetCustomerForTransactionHistoryDetail(bookingID int) (*CustomerForTrasactionHistoryDetail, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(CustomerForTrasactionHistoryDetail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetBusinessAdminInformation(userID int) (*InfoForDisbursement, error) {
	args := m.Called(userID)
	ret := args.Get(0).(InfoForDisbursement)
	return &ret, args.Error(1)
}

func (m *MockRepository) SaveDisbursement(disbursement DisbursementDetail) (int, error) {
	args := m.Called(disbursement)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) UpdateBalance(newBalance float64, userID int) error {
	args := m.Called(newBalance, userID)
	return args.Error(0)
}

func (m *MockRepository) UpdateDisbursementStatusByXenditID(newStatus int, xenditID string) error {
	args := m.Called(newStatus, xenditID)
	return args.Error(0)
}

type MockXenditService struct {
	mock.Mock
}

func (x *MockXenditService) CreateInvoice(params xendit.CreateInvoiceParams) (*xendit2.Invoice, error) {
	args := x.Called(params)
	return args.Get(0).(*xendit2.Invoice), args.Error(1)
}

func (x *MockXenditService) CreateDisbursement(params xendit.CreateDisbursementParams) (*xendit2.Disbursement, error) {
	args := x.Called(params)
	return args.Get(0).(*xendit2.Disbursement), args.Error(1)
}

func (x *MockXenditService) GetInvoice(ID string) (*xendit2.Invoice, error) {
	args := x.Called(ID)
	return args.Get(0).(*xendit2.Invoice), args.Error(1)
}

func (x *MockXenditService) GetDisbursement(ID string) (*xendit2.Disbursement, error) {
	args := x.Called(ID)
	return args.Get(0).(*xendit2.Disbursement), args.Error(1)
}

func TestService_GetBalanceDetailSuccess(t *testing.T) {
	userID := 1
	placeID := 2
	balance := BalanceDetail{
		Balance: 2500000,
	}

	latestDisbursement := DisbursementDetail{
		Date:   time.Now(),
		Amount: 500000,
		Status: 1,
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetPlaceIDByUserID", userID).Return(placeID, nil)
	mockRepo.On("GetLatestDisbursement", placeID).Return(latestDisbursement, nil)
	mockRepo.On("GetBalance", userID).Return(balance, nil)

	balanceDetailResult, err := mockService.GetBalanceDetail(userID)
	mockRepo.AssertExpectations(t)

	var balanceDetail BalanceDetail
	balanceDetail.LatestDisbursementDate = latestDisbursement.Date.String()
	balanceDetail.Balance = balance.Balance

	assert.Equal(t, &balanceDetail, balanceDetailResult)
	assert.NotNil(t, balanceDetailResult)
	assert.NoError(t, err)
}

func TestService_GetBalanceDetailWithWrongInput(t *testing.T) {
	userID := 0

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	balanceDetail, err := mockService.GetBalanceDetail(userID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, balanceDetail)
}

func TestService_GetBalanceDetailFailedCalledGetPlaceIDByUserID(t *testing.T) {
	userID := 10

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetPlaceIDByUserID", userID).Return(0, ErrInternalServerError)

	balanceDetailResult, err := mockService.GetBalanceDetail(userID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, balanceDetailResult)
}

func TestService_GetBalanceDetailFailedCalledGetLatestDisbursement(t *testing.T) {
	userID := 10
	placeID := 10
	var disbursementsDetail DisbursementDetail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetPlaceIDByUserID", userID).Return(placeID, nil)
	mockRepo.On("GetLatestDisbursement", placeID).Return(disbursementsDetail, ErrInternalServerError)

	balanceDetailResult, err := mockService.GetBalanceDetail(userID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, balanceDetailResult)
}

func TestService_GetBalanceDetailFailedCalledGetBalance(t *testing.T) {
	userID := 10
	placeID := 10
	var disbursementsDetail DisbursementDetail
	var balanceDetail BalanceDetail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetPlaceIDByUserID", userID).Return(placeID, nil)
	mockRepo.On("GetLatestDisbursement", placeID).Return(disbursementsDetail, nil)
	mockRepo.On("GetBalance", userID).Return(balanceDetail, ErrInternalServerError)

	balanceDetailResult, err := mockService.GetBalanceDetail(userID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, balanceDetailResult)
}

func TestService_GetListTransactionHistoryWithPaginationSuccess(t *testing.T) {
	// Define input and output
	listTransactionExpected := ListTransaction{
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

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	params := ListTransactionRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/testing",
		UserID: 0,
	}
	// Expectation
	mockRepo.On("GetListTransactionsHistoryWithPagination", params).Return(listTransactionExpected, nil)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listTransactionExpected, listTransactionResult)
	assert.NotNil(t, listTransactionResult)
	assert.NoError(t, err)
}

func TestService_GetListTransactionHistoryWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	listTransactionExpected := ListTransaction{
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
		Limit:  0,
		Page:   0,
		Path:   "/api/testing",
		UserID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	paramsDefault := ListTransactionRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/testing",
		UserID: 1,
	}
	// Expectation
	mockRepo.On("GetListTransactionsHistoryWithPagination", paramsDefault).Return(listTransactionExpected, nil)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listTransactionExpected, listTransactionResult)
	assert.NotNil(t, listTransactionResult)
	assert.NoError(t, err)
}

func TestService_GetListTransactionHistoryWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	params := ListTransactionRequest{
		Limit: 101,
		Page:  0,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listTransactionResult)
}

func TestService_GetListItemByIDWithPaginationError(t *testing.T) {
	listTransaction := ListTransaction{}

	params := ListTransactionRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/testing",
		UserID: 1,
	}

	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetListTransactionsHistoryWithPagination", params).Return(listTransaction, ErrInternalServerError)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listTransactionResult)
}

func TestService_GetListItemWithPaginationFailedURLIsEmpty(t *testing.T) {
	// Define input
	params := ListTransactionRequest{
		Limit: 100,
		Page:  0,
		Path:  "",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listTransactionResult)
}

func TestService_CreateDisbursement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		lastDisbursementDate := time.Date(2022, 03, 01, 0, 0, 0, 0, time.Local)

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		lastDisbursementInfo := DisbursementDetail{
			ID:       1,
			PlaceID:  1,
			Date:     lastDisbursementDate,
			XenditID: "1",
			Amount:   10000,
			Status:   0,
		}

		xenditDisbursementParams := xendit.CreateDisbursementParams{
			ID:                businessAdminInfo.ID,
			BankAccountName:   businessAdminInfo.BankAccountName,
			BankAccountNumber: businessAdminInfo.BankAccountNumber,
			Amount:            4450,
			Description:       fmt.Sprintf("Disbursement by %s", businessAdminInfo.Name),
			Email:             []string{businessAdminInfo.Email},
		}

		createXenditDisbursement := xendit2.Disbursement{ID: "1", Amount: 4450}

		disbursement := DisbursementDetail{
			PlaceID:  businessAdminInfo.PlaceID,
			Date:     time.Now(),
			XenditID: createXenditDisbursement.ID,
			Amount:   createXenditDisbursement.Amount,
			Status:   0,
		}

		expectedOutput := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: disbursement.Date,
			Amount:    4450,
			XenditID:  "1",
		}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, nil)
		mockRepo.On("GetLatestDisbursement", 1).Return(lastDisbursementInfo, nil)
		mockXendit.On("CreateDisbursement", xenditDisbursementParams).Return(&createXenditDisbursement, nil)
		mockRepo.On("SaveDisbursement", disbursement).Return(1, nil)

		resp, err := service.CreateDisbursement(1, 10000)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, &expectedOutput, resp)
	})

	t.Run("error while calling SaveDisbursement", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		lastDisbursementDate := time.Date(2022, 03, 01, 0, 0, 0, 0, time.Local)

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		lastDisbursementInfo := DisbursementDetail{
			ID:       1,
			PlaceID:  1,
			Date:     lastDisbursementDate,
			XenditID: "1",
			Amount:   10000,
			Status:   0,
		}

		xenditDisbursementParams := xendit.CreateDisbursementParams{
			ID:                businessAdminInfo.ID,
			BankAccountName:   businessAdminInfo.BankAccountName,
			BankAccountNumber: businessAdminInfo.BankAccountNumber,
			Amount:            4450,
			Description:       fmt.Sprintf("Disbursement by %s", businessAdminInfo.Name),
			Email:             []string{businessAdminInfo.Email},
		}

		createXenditDisbursement := xendit2.Disbursement{ID: "1", Amount: 4450}

		disbursement := DisbursementDetail{
			PlaceID:  businessAdminInfo.PlaceID,
			Date:     time.Now(),
			XenditID: createXenditDisbursement.ID,
			Amount:   createXenditDisbursement.Amount,
			Status:   0,
		}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, nil)
		mockRepo.On("GetLatestDisbursement", 1).Return(lastDisbursementInfo, nil)
		mockXendit.On("CreateDisbursement", xenditDisbursementParams).Return(&createXenditDisbursement, nil)
		mockRepo.On("SaveDisbursement", disbursement).Return(1, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateDisbursement(1, 10000)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("error while calling CreateDisbursement", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		lastDisbursementDate := time.Date(2022, 03, 01, 0, 0, 0, 0, time.Local)

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		lastDisbursementInfo := DisbursementDetail{
			ID:       1,
			PlaceID:  1,
			Date:     lastDisbursementDate,
			XenditID: "1",
			Amount:   10000,
			Status:   0,
		}

		xenditDisbursementParams := xendit.CreateDisbursementParams{
			ID:                businessAdminInfo.ID,
			BankAccountName:   businessAdminInfo.BankAccountName,
			BankAccountNumber: businessAdminInfo.BankAccountNumber,
			Amount:            4450,
			Description:       fmt.Sprintf("Disbursement by %s", businessAdminInfo.Name),
			Email:             []string{businessAdminInfo.Email},
		}

		createXenditDisbursement := xendit2.Disbursement{ID: "1", Amount: 4450}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, nil)
		mockRepo.On("GetLatestDisbursement", 1).Return(lastDisbursementInfo, nil)
		mockXendit.On("CreateDisbursement", xenditDisbursementParams).Return(&createXenditDisbursement, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateDisbursement(1, 10000)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("error while calling GetLatestDisbursement", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		lastDisbursementDate := time.Date(2022, 03, 01, 0, 0, 0, 0, time.Local)

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		lastDisbursementInfo := DisbursementDetail{
			ID:       1,
			PlaceID:  1,
			Date:     lastDisbursementDate,
			XenditID: "1",
			Amount:   10000,
			Status:   0,
		}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, nil)
		mockRepo.On("GetLatestDisbursement", 1).Return(lastDisbursementInfo, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateDisbursement(1, 10000)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("error while input validation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		resp, err := service.CreateDisbursement(-1, -10000)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})

	t.Run("error while calling GetBusinessAdminInformation", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, errors.Wrap(ErrInternalServerError, "test error"))

		resp, err := service.CreateDisbursement(1, 10000)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("input validation error when last disbursement is yesterday", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockXendit := new(MockXenditService)
		service := NewService(mockRepo, mockXendit)

		f := faketime.NewFaketime(2022, 04, 02, 0, 0, 0, 0, time.Local)
		defer f.Undo()
		f.Do()

		lastDisbursementDate := time.Date(2022, 04, 01, 0, 0, 0, 0, time.Local)

		businessAdminInfo := InfoForDisbursement{
			ID:                1,
			Name:              "test",
			Email:             "test@gmail.com",
			BankAccountName:   "TEST",
			BankAccountNumber: "123456789",
			PlaceID:           1,
		}

		lastDisbursementInfo := DisbursementDetail{
			ID:       1,
			PlaceID:  1,
			Date:     lastDisbursementDate,
			XenditID: "1",
			Amount:   10000,
			Status:   0,
		}

		mockRepo.On("GetBusinessAdminInformation", 1).Return(businessAdminInfo, nil)
		mockRepo.On("GetLatestDisbursement", 1).Return(lastDisbursementInfo, nil)

		resp, err := service.CreateDisbursement(1, 10000)
		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestService_DisbursementCallbackFromXendit(t *testing.T) {
	t.Run("success status completed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "COMPLETED",
		}

		balanceDetail := BalanceDetail{Balance: 10000, LatestDisbursementDate: time.Now().String()}
		mockRepo.On("GetBalance", 1).Return(balanceDetail, nil)
		mockRepo.On("UpdateBalance", 0.0, 1).Return(nil)
		mockRepo.On("UpdateDisbursementStatusByXenditID", 1, "test").Return(nil)

		err := service.DisbursementCallbackFromXendit(params)
		assert.Nil(t, err)
	})

	t.Run("failed status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "TEST",
		}

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})

	t.Run("failed calling update disbursement status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "COMPLETED",
		}

		balanceDetail := BalanceDetail{Balance: 10000, LatestDisbursementDate: time.Now().String()}
		mockRepo.On("GetBalance", 1).Return(balanceDetail, nil)
		mockRepo.On("UpdateBalance", 0.0, 1).Return(nil)
		mockRepo.On("UpdateDisbursementStatusByXenditID", 1, "test").Return(errors.Wrap(ErrInternalServerError, "test error"))

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed calling update disbursement status on failed callback case", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "FAILED",
		}

		mockRepo.On("UpdateDisbursementStatusByXenditID", 2, "test").Return(errors.Wrap(ErrInternalServerError, "test error"))

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed calling update balance", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "COMPLETED",
		}

		balanceDetail := BalanceDetail{Balance: 10000, LatestDisbursementDate: time.Now().String()}
		mockRepo.On("GetBalance", 1).Return(balanceDetail, nil)
		mockRepo.On("UpdateBalance", 0.0, 1).Return(errors.Wrap(ErrInternalServerError, "test error"))

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("failed when calling get balance", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "COMPLETED",
		}

		balanceDetail := BalanceDetail{Balance: 10000, LatestDisbursementDate: time.Now().String()}
		mockRepo.On("GetBalance", 1).Return(balanceDetail, errors.Wrap(ErrInternalServerError, "test error"))

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Error(t, ErrInternalServerError, errors.Cause(err))
	})

	t.Run("success status failed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "1",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "FAILED",
		}

		mockRepo.On("UpdateDisbursementStatusByXenditID", 2, "test").Return(nil)

		err := service.DisbursementCallbackFromXendit(params)
		assert.Nil(t, err)
	})

	t.Run("failed parse user id", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		params := DisbursementCallback{
			ID:                      "test",
			ExternalID:              "test",
			Amount:                  10000,
			BankCode:                "BCA",
			AccountHolderName:       "TEST",
			DisbursementDescription: "test",
			FailureCode:             "",
			Status:                  "FAILED",
		}

		err := service.DisbursementCallbackFromXendit(params)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	})
}

func TestService_GetTransactionHistoryDetailWithWrongInput(t *testing.T) {
	bookingID := 0

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	balanceDetail, err := mockService.GetTransactionHistoryDetail(bookingID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, balanceDetail)
}

func TestService_GetTransactionHistoryDetailSuccess(t *testing.T) {
	bookingID := 1
	itemsWrapper := ItemsWrapper{
		Items: []ItemDetail{
			{
				Name:  "ini_nama_item_1",
				Qty:   25,
				Price: 10000,
			},
			{
				Name:  "ini_nama_item_2",
				Qty:   5,
				Price: 20000,
			},
		},
	}

	customer := CustomerForTrasactionHistoryDetail{
		CustomerName:  "ini_customer_name",
		CustomerImage: "ini_customer_image",
	}

	transactionHistoryDetail := TransactionHistoryDetail{
		Date:           "27 Oktober 2021",
		StartTime:      "08:00",
		EndTime:        "09:00",
		Capacity:       20,
		TotalPriceItem: 25000,
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemsWrapper", bookingID).Return(itemsWrapper, nil)
	mockRepo.On("GetCustomerForTransactionHistoryDetail", bookingID).Return(customer, nil)
	mockRepo.On("GetTransactionHistoryDetail", bookingID).Return(transactionHistoryDetail, nil)

	transactionHistoryDetailResult, err := mockService.GetTransactionHistoryDetail(bookingID)
	mockRepo.AssertExpectations(t)

	transactionHistoryDetail.CustomerName = customer.CustomerName
	transactionHistoryDetail.CustomerImage = customer.CustomerImage
	transactionHistoryDetail.Items = itemsWrapper.Items

	assert.Equal(t, &transactionHistoryDetail, transactionHistoryDetailResult)
	assert.NotNil(t, transactionHistoryDetailResult)
	assert.NoError(t, err)
}

func TestService_GetTransactionHistoryDetailFailedCalledGetItemsWrapper(t *testing.T) {
	bookingID := 1
	var itemsWrapper ItemsWrapper

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemsWrapper", bookingID).Return(itemsWrapper, ErrInternalServerError)

	transactionHistoryDetailResult, err := mockService.GetTransactionHistoryDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, transactionHistoryDetailResult)
}

func TestService_GetTransactionHistoryDetailFailedCalledGetCustomerForTransactionHistoryDetail(t *testing.T) {
	bookingID := 1
	var itemsWrapper ItemsWrapper
	var customer CustomerForTrasactionHistoryDetail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemsWrapper", bookingID).Return(itemsWrapper, nil)
	mockRepo.On("GetCustomerForTransactionHistoryDetail", bookingID).Return(customer, ErrInternalServerError)

	transactionHistoryDetailResult, err := mockService.GetTransactionHistoryDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, transactionHistoryDetailResult)
}

func TestService_GetTransactionHistoryDetailFailedCalledGetTransactionHistoryDetail(t *testing.T) {
	bookingID := 1
	var itemsWrapper ItemsWrapper
	var customer CustomerForTrasactionHistoryDetail
	var transactionHistoryDetail TransactionHistoryDetail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemsWrapper", bookingID).Return(itemsWrapper, nil)
	mockRepo.On("GetCustomerForTransactionHistoryDetail", bookingID).Return(customer, nil)
	mockRepo.On("GetTransactionHistoryDetail", bookingID).Return(transactionHistoryDetail, ErrInternalServerError)

	transactionHistoryDetailResult, err := mockService.GetTransactionHistoryDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, transactionHistoryDetailResult)
}
