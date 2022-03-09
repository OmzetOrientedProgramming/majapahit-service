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

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockService) GetItemByID(placeID int, itemID int ) (*Item, error) {
	args := m.Called(placeID, itemID)
	item := args.Get(0).(*Item)
	return item, args.Error(1)
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
		Limit: 10,
		Page:  1,
		Path:  "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          	1,
				Name:        	"test",
				Image:     		"test",
				Description:	"test",
				Price:    		10000,
			},
			{
				ID:          	2,
				Name:        	"test",
				Image:     		"test",
				Description:	"test",
				Price:    		10000,
			},
		},
		TotalCount: 10,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstUrl:    fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		LastUrl:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		NextUrl:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousUrl: fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":     listItem.Items,
			"pagination": pagination,
		},
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Status: http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect place id",
		},

	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetListItemWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Limit: 110,
		Page:  1,
		Path:  "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, errorFromService)


	// Tes
	assert.NoError(t, h.GetListItemWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetListItemWithPagination(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Limit: 110,
		Page:  1,
		Path:  "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listItem ListItem
	var pagination util.Pagination
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, internalServerError)

	// Tes
	assert.NoError(t, h.GetListItemWithPagination(ctx))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Limit: 0,
		Page:  0,
		Path:  "/api/v1/place/1/catalog",
		PlaceID: 1,
	}

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	listItem := ListItem{
		Items: []Item{
			{
				ID:          	1,
				Name:        	"test",
				Image:     		"test",
				Description:	"test",
				Price:    		10000,
			},
			{
				ID:          	2,
				Name:        	"test",
				Image:     		"test",
				Description:	"test",
				Price:    		10000,
			},
		},
		TotalCount: 10,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstUrl:    fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		LastUrl:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		NextUrl:     fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousUrl: fmt.Sprintf("%s/api/v1/place/1/catalog?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":     listItem.Items,
			"pagination": pagination,
		},
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItemWithPagination", params).Return(&listItem, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetListItemWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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

	item := Item {
			ID:          	1,
			Name:        	"test",
			Image:     		"test",
			Description:	"test",
			Price:    		10000,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"item":     item,
		},
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetItemByID", placeID, itemID).Return(&item, nil)

	// Tes
	if assert.NoError(t, h.GetItemByID(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Status: http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect place id",
		},

	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetItemByID(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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
		Status: http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"incorrect item id",
		},

	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Tes
	assert.NoError(t, h.GetItemByID(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
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

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	var item Item
	mockService.On("GetItemByID", placeID, itemID).Return(&item, internalServerError)

	// Tes
	assert.NoError(t, h.GetItemByID(ctx))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
}