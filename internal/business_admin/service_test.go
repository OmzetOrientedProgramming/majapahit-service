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
	balanceDetail.Balance = balance.Balance - float64(latestDisbursement.Amount)

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
