package businessadmin

import (
	"testing"

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

func TestService_GetBalanceDetailSuccess(t *testing.T) {
	userID := 1
	placeID := 2
	balance := BalanceDetail{
		Balance: 2500000,
	}

	latestDisbursement := DisbursementDetail{
		Date:   "27 Januari 2022",
		Amount: 500000,
		Status: 1,
	}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetPlaceIDByUserID", userID).Return(placeID, nil)
	mockRepo.On("GetLatestDisbursement", placeID).Return(latestDisbursement, nil)
	mockRepo.On("GetBalance", userID).Return(balance, nil)

	balanceDetailResult, err := mockService.GetBalanceDetail(userID)
	mockRepo.AssertExpectations(t)

	var balanceDetail BalanceDetail
	balanceDetail.LatestDisbursementDate = latestDisbursement.Date
	balanceDetail.Balance = balance.Balance

	assert.Equal(t, &balanceDetail, balanceDetailResult)
	assert.NotNil(t, balanceDetailResult)
	assert.NoError(t, err)
}

func TestService_GetBalanceDetailWithWrongInput(t *testing.T) {
	userID := 0

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	balanceDetail, err := mockService.GetBalanceDetail(userID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, balanceDetail)
}

func TestService_GetBalanceDetailFailedCalledGetPlaceIDByUserID(t *testing.T) {
	userID := 10

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

	// Test
	listTransactionResult, _, err := mockService.GetListTransactionsHistoryWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listTransactionResult)
}

func TestService_GetTransactionHistoryDetailWithWrongInput(t *testing.T) {
	bookingID := 0

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

	mockRepo.On("GetItemsWrapper", bookingID).Return(itemsWrapper, nil)
	mockRepo.On("GetCustomerForTransactionHistoryDetail", bookingID).Return(customer, nil)
	mockRepo.On("GetTransactionHistoryDetail", bookingID).Return(transactionHistoryDetail, ErrInternalServerError)

	transactionHistoryDetailResult, err := mockService.GetTransactionHistoryDetail(bookingID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, transactionHistoryDetailResult)
}
