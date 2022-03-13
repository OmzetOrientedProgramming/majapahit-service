package place

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	args := m.Called(placeId)
	placeDetail := args.Get(0).(*PlaceDetail)
	return placeDetail, args.Error(1)
}

func TestHandler_GetPlaceDetailSuccess(t *testing.T) {
	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:placeId")
	c.SetParamNames("placeId")
	c.SetParamValues("1")

	// Setting up service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setting up Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Setting up input and output
	placeId := 1

	placeDetail := PlaceDetail{
		ID:            1,
		Name:          "test_name_place",
		Image:         "test_image_place",
		Distance:      200,
		Address:       "test_address_place",
		Description:   "test_description_place",
		OpenHour:      "08:00",
		CloseHour:     "16:00",
		AverageRating: 3.50,
		ReviewCount:   30,
		Reviews: []UserReview{
			{
				User:    "test_user_1",
				Rating:  4.50,
				Content: "test_review_content_1",
			},
			{
				User:    "test_user_2",
				Rating:  5,
				Content: "test_review_content_2",
			},
		},
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    placeDetail,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetPlaceDetail", placeId).Return(&placeDetail, nil)

	// Test Fields
	if assert.NoError(t, h.GetPlaceDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func (m *MockService) GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error) {
	args := m.Called(params)
	placeList := args.Get(0).(*PlacesList)
	pagination := args.Get(1).(util.Pagination)
	return placeList, &pagination, args.Error(2)
}

func TestHandler_GetPlacesListWithPaginationWithParams(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/v1/place",
	}

	placeList := PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test name",
				Description: "test description",
				Address:     "test address",
				Distance:    10,
				Rating:      4.5,
				ReviewCount: 20,
			},
			{
				ID:          2,
				Name:        "test name 2",
				Description: "test description 2",
				Address:     "test address 2",
				Distance:    11,
				Rating:      2.0,
				ReviewCount: 100,
			},
		},
		TotalCount: 2,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"places":     placeList.Places,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetPlaceListWithPagination", params).Return(&placeList, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetPlacesListWithPagination(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetPlacesListWithPaginationWithParamsError(t *testing.T) {
	// Setup echo
	e := echo.New()
	q := make(url.Values)
	q.Set("limit", "1001")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := PlacesListRequest{
		Limit: 1001,
		Page:  1,
		Path:  "/api/v1/place",
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var placeList PlacesList
	var pagination util.Pagination
	mockService.On("GetPlaceListWithPagination", params).Return(&placeList, pagination, errorFromService)

	// Tes
	assert.NoError(t, h.GetPlacesListWithPagination(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetPlacesListWithPaginationWithInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	q := make(url.Values)
	q.Set("limit", "1001")
	q.Set("page", "1")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := PlacesListRequest{
		Limit: 1001,
		Page:  1,
		Path:  "/api/v1/place",
	}

	errorFromService := errors.Wrap(ErrInternalServerError, "test error")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var placeList PlacesList
	var pagination util.Pagination
	mockService.On("GetPlaceListWithPagination", params).Return(&placeList, pagination, errorFromService)

	// Tes
	assert.NoError(t, h.GetPlacesListWithPagination(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetPlacesListWithPaginationWithValidationErrorLimitPageNotInt(t *testing.T) {
	// Setup echo
	e := echo.New()
	q := make(url.Values)
	q.Set("limit", "testerror")
	q.Set("page", "testerror")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup service
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
	assert.NoError(t, h.GetPlacesListWithPagination(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetPlacesListWithPaginationWithoutParams(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Define input and output
	params := PlacesListRequest{
		Limit: 0,
		Page:  0,
		Path:  "/api/v1/place",
	}

	placeList := PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test name",
				Description: "test description",
				Address:     "test address",
				Distance:    10,
				Rating:      4.5,
				ReviewCount: 20,
			},
			{
				ID:          2,
				Name:        "test name 2",
				Description: "test description 2",
				Address:     "test address 2",
				Distance:    11,
				Rating:      2.0,
				ReviewCount: 100,
			},
		},
		TotalCount: 2,
	}

	pagination := util.Pagination{
		Limit:       10,
		Page:        1,
		FirstURL:    fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		LastURL:     fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		NextURL:     fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		PreviousURL: fmt.Sprintf("%s/api/v1/place?limit=10&page=1", os.Getenv("BASE_URL")),
		TotalPage:   1,
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"places":     placeList.Places,
			"pagination": pagination,
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetPlaceListWithPagination", params).Return(&placeList, pagination, nil)

	// Tes
	if assert.NoError(t, h.GetPlacesListWithPagination(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}
