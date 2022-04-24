package item

import (
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
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetListItemWithPagination(params ListItemRequest) (*ListItem, *util.Pagination, error) {
	args := m.Called(params)
	listItem := args.Get(0).(*ListItem)
	pagination := args.Get(1).(util.Pagination)
	return listItem, &pagination, args.Error(2)
}

func (m *MockService) GetItemByID(placeID int, itemID int) (*Item, error) {
	args := m.Called(placeID, itemID)
	item := args.Get(0).(*Item)
	return item, args.Error(1)
}

func (m *MockService) DeleteItemAdminByID(itemID int) error {
	args := m.Called(itemID)
	return args.Error(0)
}

func TestHandler_GetListItemWithPaginationSuccess(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
		},
		TotalCount: 10,
		PlaceInfo: []PlaceInfo{
			{
				Name:  "test",
				Image: "test",
			},
		},
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":      listItem.Items,
			"info":       listItem.PlaceInfo,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListItemWithPaginationPlaceIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("test")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect place id",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	util.ErrorHandler(h.GetListItemWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemWithPaginationLimitError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListItemRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, errorFromService)
	util.ErrorHandler(h.GetListItemWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemWithPaginationLimitAndPageAreNotInt(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	q.Set("limit", "asd")
	q.Set("page", "asd")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

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
	util.ErrorHandler(h.GetListItemWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	q.Set("limit", "110")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListItemRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, internalServerError)

	// Tes
	util.ErrorHandler(h.GetListItemWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemWithPaginationWithLimitAndPageAreEmpty(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
		},
		TotalCount: 10,
		PlaceInfo: []PlaceInfo{
			{
				Name:  "test",
				Image: "test",
			},
		},
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":      listItem.Items,
			"info":       listItem.PlaceInfo,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetItemByID(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode()+"/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog/:itemID")
	ctx.SetParamNames("placeID", "itemID")
	ctx.SetParamValues("10", "1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	placeID := 10
	itemID := 1

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	item := Item{
		ID:          1,
		Name:        "test",
		Image:       "test",
		Description: "test",
		Price:       10000,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"item": item,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetItemByID", placeID, itemID).Return(&item, nil)

	// Tes
	if assert.NoError(t, h.GetItemByID(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetItemByIDPlaceIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode()+"/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog/:itemID")
	ctx.SetParamNames("placeID", "itemID")
	ctx.SetParamValues("test", "1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect place id",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	util.ErrorHandler(h.GetItemByID(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetItemByIDItemIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode()+"/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog/:itemID")
	ctx.SetParamNames("placeID", "itemID")
	ctx.SetParamValues("1", "test")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect item id",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	util.ErrorHandler(h.GetItemByID(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetItemByIDInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode()+"/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog/:itemID")
	ctx.SetParamNames("placeID", "itemID")
	ctx.SetParamValues("10", "1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	placeID := 10
	itemID := 1

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var item Item
	mockService.On("GetItemByID", placeID, itemID).Return(&item, internalServerError)

	// Tes
	util.ErrorHandler(h.GetItemByID(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemAdminWithPaginationSuccess(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/list-items",
		UserID:  1,
		PlaceID: 0,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
		},
		TotalCount: 10,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":      listItem.Items,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemAdminWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListItemAdminWithPaginationParseUserDataError(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/list-items",
		UserID:  1,
		PlaceID: 0,
	}

	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	util.ErrorHandler(h.GetListItemAdminWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestHandler_GetListItemAdminWithPaginationLimitError(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListItemRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/list-items",
		UserID:  1,
		PlaceID: 0,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, errorFromService)

	// Tes
	util.ErrorHandler(h.GetListItemAdminWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemAdminWithPaginationLimitAndPageAreNotInt(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
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
	util.ErrorHandler(h.GetListItemAdminWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemAdminInternalServerError(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListItemRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/list-items",
		UserID:  1,
		PlaceID: 0,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, internalServerError)

	// Tes
	util.ErrorHandler(h.GetListItemAdminWithPagination(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemAdminWithPaginationWithLimitAndPageAreEmpty(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodGet, "/list-items?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("userFromDatabase", &userModel)
	ctx.Set("userFromFirebase", &userData)

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/business-admin/business-profile/list-items",
		UserID:  1,
		PlaceID: 0,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test",
				Image:       "test",
				Description: "test",
				Price:       10000,
			},
		},
		TotalCount: 10,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/business-admin/business-profile/list-items?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":      listItem.Items,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemAdminWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_DeleteItemAdminByID(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/business-admin/business-profile/list-items", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:itemID")
	ctx.SetParamNames("itemID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	itemID := 1

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("DeleteItemAdminByID", itemID).Return(nil)

	// Tes
	if assert.NoError(t, h.DeleteItemAdminByID(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_DeleteItemAdminByIDItemIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/list-items/:itemID")
	ctx.SetParamNames("itemID")
	ctx.SetParamValues("test")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect item id",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Tes
	util.ErrorHandler(h.DeleteItemAdminByID(ctx), ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_DeleteItemAdminByIDInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/api/v1/business-admin/business-profile/list-items/:itemID")
	ctx.SetParamNames("itemID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	itemID := 1

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("DeleteItemAdminByID", itemID).Return(internalServerError)

	// Tes
	util.ErrorHandler(h.DeleteItemAdminByID(ctx), ctx)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}