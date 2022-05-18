package businessadmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetBalanceDetail(userID int) (*BalanceDetail, error) {
	args := m.Called(userID)
	ret := args.Get(0).(*BalanceDetail)
	return ret, args.Error(1)
}

func (m *MockService) GetListTransactionsHistoryWithPagination(params ListTransactionRequest) (*ListTransaction, *util.Pagination, error) {
	args := m.Called(params)
	listItem := args.Get(0).(*ListTransaction)
	pagination := args.Get(1).(util.Pagination)
	return listItem, &pagination, args.Error(2)
}

func (m *MockService) CreateDisbursement(ID int, amount float64) (*CreateDisbursementResponse, error) {
	args := m.Called(ID, amount)
	return args.Get(0).(*CreateDisbursementResponse), args.Error(1)
}

func (m *MockService) DisbursementCallbackFromXendit(params DisbursementCallback) error {
	args := m.Called(params)
	return args.Error(0)
}

func (m *MockService) GetTransactionHistoryDetail(bookingID int) (*TransactionHistoryDetail, error) {
	args := m.Called(bookingID)
	ret := args.Get(0).(*TransactionHistoryDetail)
	return ret, args.Error(1)
}

func (m *MockService) GetPlaceDetail(userID int) (*PlaceDetail, error) {
	args := m.Called(userID)
	ret := args.Get(0).(*PlaceDetail)
	return ret, args.Error(1)
}

func (m *MockService) GetListReviewAndRatingWithPagination(userID int, params ListReviewRequest) (*place.ListReview, *util.Pagination, error) {
	args := m.Called(userID, params)
	listReview := args.Get(0).(*place.ListReview)
	pagination := args.Get(1).(util.Pagination)
	return listReview, &pagination, args.Error(2)
}

func TestHandler_GetBalanceDetailSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	balanceDetail := BalanceDetail{
		LatestDisbursementDate: "27 Januari 2022",
		Balance:                2500000,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    balanceDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, nil)

	// Tes
	if assert.NoError(t, h.GetBalanceDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetBalanceDetailParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Excpectation
	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, nil)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetBalanceDetailInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetBalanceDetailBadRequestFromService(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/balance")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInputValidationError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"test",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var balanceDetail BalanceDetail
	mockService.On("GetBalanceDetail", userModel.ID).Return(&balanceDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetBalanceDetail(c), c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListTransactionHistoryWithPaginationSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/transaction-history"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListTransactionRequest{
		Limit:  10,
		Page:   1,
		Path:   "/api/v1/business-admin/transaction-history",
		UserID: userModel.ID,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listTransaction := ListTransaction{
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

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/transaction-history?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/transaction-history?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/transaction-history?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/transaction-history?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"transactions": listTransaction.Transactions,
			"pagination":   pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListTransactionsHistoryWithPagination", params).Return(&listTransaction, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListTransactionsHistoryWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListTransactionsHistoryWithPaginationStateAndLimitAndPageAreNotInt(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "asd")
	q.Set("page", "asd")
	req := httptest.NewRequest(http.MethodGet, "/transaction-history?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"limit should be positive integer",
			"page should be positive integer",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	util.ErrorHandler(h.GetListTransactionsHistoryWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListTransactionsHistoryWithPaginationLimitError(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/transaction-history?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListTransactionRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/transaction-history",
		UserID: 1,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listTransaction ListTransaction
	var pagination util.Pagination
	mockService.On("GetListTransactionsHistoryWithPagination", params).Return(&listTransaction, pagination, errorFromService)

	// Tes
	util.ErrorHandler(h.GetListTransactionsHistoryWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListCustomerBookingWithPaginationInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/transaction-history?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListTransactionRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/transaction-history",
		UserID: 1,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listTransaction ListTransaction
	var pagination util.Pagination
	mockService.On("GetListTransactionsHistoryWithPagination", params).Return(&listTransaction, pagination, internalServerError)

	// Tes
	util.ErrorHandler(h.GetListTransactionsHistoryWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListTransactionsHistoryWithPaginationParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/transaction-history?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListTransactionRequest{
		Limit:  110,
		Page:   1,
		Path:   "/api/v1/business-admin/transaction-history",
		UserID: 1,
	}

	// Excpectation
	var listTransaction ListTransaction
	var pagination util.Pagination
	mockService.On("GetListTransactionsHistoryWithPagination", params).Return(&listTransaction, pagination, nil)

	// Tes
	util.ErrorHandler(h.GetListTransactionsHistoryWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_CreateDisbursement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": 10000})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  200,
			Message: "success",
			Data:    respData,
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, nil)
		err := h.CreateDisbursement(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed internal server error from service", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": 10000})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  500,
			Message: "internal server error",
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, errors.Wrap(ErrInternalServerError, "tes error"))
		err := h.CreateDisbursement(ctx)

		util.ErrorHandler(err, ctx)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed input validation error from service", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": 10000})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  400,
			Message: "input validation error",
			Errors:  []string{"test error"},
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, errors.Wrap(ErrInputValidationError, "test error"))
		err := h.CreateDisbursement(ctx)

		util.ErrorHandler(err, ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed bind body", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": 10000})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  400,
			Message: "input validation error",
			Errors:  []string{"invalid body"},
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, nil)
		err := h.CreateDisbursement(ctx)

		util.ErrorHandler(err, ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed amount", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "password",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": "test"})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  400,
			Message: "input validation error",
			Errors:  []string{"invalid amount"},
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, nil)
		err := h.CreateDisbursement(ctx)

		util.ErrorHandler(err, ctx)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("success", func(t *testing.T) {
		// setup echo
		e := echo.New()
		userData := firebaseauth.UserDataFromToken{
			Kind: "",
			Users: []firebaseauth.User{
				{
					LocalID: "1",
					ProviderUserInfo: []firebaseauth.ProviderUserInfo{
						{
							ProviderID:  "phone",
							RawID:       "",
							PhoneNumber: "",
							FederatedID: "",
							Email:       "",
						},
					},
					LastLoginAt:       "",
					CreatedAt:         "",
					PhoneNumber:       "",
					LastRefreshAt:     time.Time{},
					Email:             "",
					EmailVerified:     false,
					PasswordHash:      "",
					PasswordUpdatedAt: 0,
					ValidSince:        "",
					Disabled:          false,
				},
			},
		}

		userModel := user.Model{
			ID:              1,
			PhoneNumber:     "",
			Name:            "",
			Status:          0,
			FirebaseLocalID: "",
			Email:           "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}

		// input
		payload, _ := json.Marshal(map[string]interface{}{"amount": 10000})

		// output
		respData := CreateDisbursementResponse{
			ID:        1,
			CreatedAt: time.Now(),
			Amount:    10000,
			XenditID:  "test xendit id",
		}

		expectedOutput := util.APIResponse{
			Status:  403,
			Message: "user does not have access to this endpoint",
			Errors: []string{
				"user is not business admin",
			},
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.Set("userFromDatabase", &userModel)
		ctx.Set("userFromFirebase", &userData)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("CreateDisbursement", 1, 10000.0).Return(&respData, nil)
		err := h.CreateDisbursement(ctx)

		util.ErrorHandler(err, ctx)
		assert.Equal(t, http.StatusForbidden, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_XenditDisbursementCallback(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// setup echo
		e := echo.New()

		// input
		input := DisbursementCallback{
			ID:                      "1",
			ExternalID:              "1",
			Amount:                  1000,
			BankCode:                "BCA",
			AccountHolderName:       "test",
			DisbursementDescription: "test description",
			FailureCode:             "TEST",
			Status:                  "COMPLETED",
		}
		payload, _ := json.Marshal(input)

		// input
		expectedOutput := util.APIResponse{
			Status:  201,
			Message: "success",
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("DisbursementCallbackFromXendit", input).Return(nil)
		err := h.XenditDisbursementCallback(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed bind param", func(t *testing.T) {
		// setup echo
		e := echo.New()

		// input
		input := map[string]interface{}{
			"id": 1,
		}
		payload, _ := json.Marshal(input)

		// input
		expectedOutput := util.APIResponse{
			Status:  500,
			Message: "internal server error",
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("DisbursementCallbackFromXendit", input).Return(nil)
		err := h.XenditDisbursementCallback(ctx)

		util.ErrorHandler(err, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed from service", func(t *testing.T) {
		// setup echo
		e := echo.New()

		// input
		input := DisbursementCallback{
			ID:                      "1",
			ExternalID:              "1",
			Amount:                  1000,
			BankCode:                "BCA",
			AccountHolderName:       "test",
			DisbursementDescription: "test description",
			FailureCode:             "TEST",
			Status:                  "COMPLETED",
		}
		payload, _ := json.Marshal(input)

		// input
		expectedOutput := util.APIResponse{
			Status:  500,
			Message: "internal server error",
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("DisbursementCallbackFromXendit", input).Return(errors.Wrap(ErrInternalServerError, "test error"))
		err := h.XenditDisbursementCallback(ctx)

		util.ErrorHandler(err, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})

	t.Run("failed from service input validation error", func(t *testing.T) {
		// setup echo
		e := echo.New()

		// input
		input := DisbursementCallback{
			ID:                      "1",
			ExternalID:              "1",
			Amount:                  1000,
			BankCode:                "BCA",
			AccountHolderName:       "test",
			DisbursementDescription: "test description",
			FailureCode:             "TEST",
			Status:                  "COMPLETED",
		}
		payload, _ := json.Marshal(input)

		// input
		expectedOutput := util.APIResponse{
			Status:  400,
			Message: "input validation error",
			Errors: []string{
				"test error",
			},
		}

		expectedOutputJSON, _ := json.Marshal(expectedOutput)

		// import "net/url"
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		mockService := new(MockService)
		h := NewHandler(mockService)

		mockService.On("DisbursementCallbackFromXendit", input).Return(errors.Wrap(ErrInputValidationError, "test error"))
		err := h.XenditDisbursementCallback(ctx)

		util.ErrorHandler(err, ctx)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, string(expectedOutputJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	})
}

func TestHandler_GetTransactionHistoryDetailSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/transaction-history/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("1")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	bookingID := 1

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	transactionHistoryDetail := TransactionHistoryDetail{
		Date:           "27 Oktober 2021",
		StartTime:      "08:00",
		EndTime:        "09:00",
		Capacity:       20,
		TotalPriceItem: 25000,
		CustomerName:   "ini_customer_name",
		CustomerImage:  "ini_customer_image",
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

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    transactionHistoryDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetTransactionHistoryDetail", bookingID).Return(&transactionHistoryDetail, nil)

	// Tes
	if assert.NoError(t, h.GetTransactionHistoryDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetTransactionHistoryDetailParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/transaction-history/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("1")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Excpectation
	var transactionHistoryDetail TransactionHistoryDetail
	mockService.On("GetTransactionHistoryDetail", userModel.ID).Return(&transactionHistoryDetail, nil)

	// Tes
	util.ErrorHandler(h.GetTransactionHistoryDetail(c), c)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetTransactionHistoryDetailWithBookingIDString(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/transaction-history/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("satu")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Excpectation
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"bookingID must be number",
		},
	}
	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	response := h.GetTransactionHistoryDetail(c)
	util.ErrorHandler(response, c)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetTransactionHistoryDetailWithBookingIDBelowOne(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/transaction-history/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("0")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	bookingID := 0

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"bookingID must be above 0"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}
	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var transactionHistoryDetail TransactionHistoryDetail
	mockService.On("GetTransactionHistoryDetail", bookingID).Return(&transactionHistoryDetail, errorFromService)

	response := h.GetTransactionHistoryDetail(c)
	util.ErrorHandler(response, c)

	// Test
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetTransactionHistoryDetailInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/transaction-history/:bookingID")
	c.SetParamNames("bookingID")
	c.SetParamValues("10")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input and output
	bookingID := 10

	errorFromService := errors.Wrap(ErrInternalServerError, "test error")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var transactionHistoryDetail TransactionHistoryDetail
	mockService.On("GetTransactionHistoryDetail", bookingID).Return(&transactionHistoryDetail, errorFromService)

	response := h.GetTransactionHistoryDetail(c)
	util.ErrorHandler(response, c)

	// Tes
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetPlaceDetailSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/business-profile/detail")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	placeDetail := PlaceDetail{
		ID:            1,
		Name:          "test_name_place",
		Image:         "test_image_place",
		Address:       "test_address_place",
		Description:   "test_description_place",
		OpenHour:      "08:00",
		CloseHour:     "16:00",
		AverageRating: 3.50,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    placeDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetPlaceDetail", userModel.ID).Return(&placeDetail, nil)

	// Tes
	if assert.NoError(t, h.GetPlaceDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetPlaceDetailParseUserDataError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/business-profile/detail")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Excpectation
	var balanceDetail BalanceDetail
	mockService.On("GetPlaceDetail", userModel.ID).Return(&balanceDetail, nil)

	// Tes
	util.ErrorHandler(h.GetPlaceDetail(c), c)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetPlaceDetailInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/business-profile/detail")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var placeDetail PlaceDetail
	mockService.On("GetPlaceDetail", userModel.ID).Return(&placeDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetPlaceDetail(c), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetPlaceDetailBadRequestFromService(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/business-admin/business-profile/detail")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	c.Set("userFromDatabase", &userModel)
	c.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	internalServerError := errors.Wrap(ErrInputValidationError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"test",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var placeDetail PlaceDetail
	mockService.On("GetPlaceDetail", userModel.ID).Return(&placeDetail, internalServerError)

	// Tes
	util.ErrorHandler(h.GetPlaceDetail(c), c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListReviewAndRatingWithPaginationSuccess(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/review")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)


	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/review",
	}

	t.Setenv("BASE_URL", "localhost:8080")

	listReview := place.ListReview{
		Reviews: []place.Review{
			{
				ID:      2,
				Name:    "test 2",
				Content: "test 2",
				Rating:  2,
				Date:    "test 2",
			},
			{
				ID:      1,
				Name:    "test 1",
				Content: "test 1",
				Rating:  1,
				Date:    "test 1",
			},
		},
		TotalCount: 2,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"reviews":      listReview.Reviews,
			"pagination":   pagination,
			"total_review": listReview.TotalCount,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	mockService.On("GetListReviewAndRatingWithPagination", userModel.ID, params).Return(&listReview, pagination, nil)

	if assert.NoError(t, h.GetListReviewAndRatingWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListReviewAndRatingWithPaginationParseUserDataError(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/review")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "phone",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)


	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/review",
	}

	t.Setenv("BASE_URL", "localhost:8080")

	var listReview place.ListReview
	var pagination util.Pagination
	mockService.On("GetListReviewAndRatingWithPagination", userModel.ID, params).Return(&listReview, pagination, nil)
	
	// Tes
	util.ErrorHandler(h.GetListReviewAndRatingWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetListReviewAndRatingWithPaginationLimitError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "101")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/review")

	params := ListReviewRequest{
		Limit:   101,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/review",
	}

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listReview place.ListReview
	var pagination util.Pagination
	mockService.On("GetListReviewAndRatingWithPagination", userModel.ID, params).Return(&listReview, pagination, errorFromService)
	util.ErrorHandler(h.GetListReviewAndRatingWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListReviewAndRatingWithPaginationInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/review")

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/review",
	}

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listReview place.ListReview
	var pagination util.Pagination
	mockService.On("GetListReviewAndRatingWithPagination", userModel.ID, params).Return(&listReview, pagination, internalServerError)
	util.ErrorHandler(h.GetListReviewAndRatingWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListReviewAndRatingWithPaginationQueryParamEmpty(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("limit", "")
	q.Set("page", "")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/review")

	userData := firebaseauth.UserDataFromToken{
		Kind: "",
		Users: []firebaseauth.User{
			{
				LocalID: "1",
				ProviderUserInfo: []firebaseauth.ProviderUserInfo{
					{
						ProviderID:  "password",
						RawID:       "",
						PhoneNumber: "",
						FederatedID: "",
						Email:       "",
					},
				},
				LastLoginAt:       "",
				CreatedAt:         "",
				PhoneNumber:       "",
				LastRefreshAt:     time.Time{},
				Email:             "",
				EmailVerified:     false,
				PasswordHash:      "",
				PasswordUpdatedAt: 0,
				ValidSince:        "",
				Disabled:          false,
			},
		},
	}

	userModel := user.Model{
		ID:              1,
		PhoneNumber:     "",
		Name:            "",
		Status:          0,
		FirebaseLocalID: "",
		Email:           "",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	paramsDefault := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/review",
	}

	t.Setenv("BASE_URL", "localhost:8080")

	listReview := place.ListReview{
		Reviews: []place.Review{
			{
				ID:      2,
				Name:    "test 2",
				Content: "test 2",
				Rating:  2,
				Date:    "test 2",
			},
			{
				ID:      1,
				Name:    "test 1",
				Content: "test 1",
				Rating:  1,
				Date:    "test 1",
			},
		},
		TotalCount: 2,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place/1/review?limit=10&page=1&latest=true&rating=true", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"reviews":      listReview.Reviews,
			"pagination":   pagination,
			"total_review": listReview.TotalCount,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	mockService.On("GetListReviewAndRatingWithPagination", userModel.ID, paramsDefault).Return(&listReview, pagination, nil)

	if assert.NoError(t, h.GetListReviewAndRatingWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}
