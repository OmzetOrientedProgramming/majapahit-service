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

func (m *MockService) GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error) {
	args := m.Called(params)
	placeList := args.Get(0).(*PlacesList)
	pagination := args.Get(1).(util.Pagination)
	return placeList, &pagination, args.Error(2)
}

func (m *MockService) GetDetail(placeID int) (*Detail, error) {
	args := m.Called(placeID)
	placeDetail := args.Get(0).(*Detail)
	return placeDetail, args.Error(1)
}

func (m *MockService) GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, *util.Pagination, error) {
	args := m.Called(params)
	listReview := args.Get(0).(*ListReview)
	pagination := args.Get(1).(util.Pagination)
	return listReview, &pagination, args.Error(2)
}

func TestHandler_GetDetailSuccess(t *testing.T) {
	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:placeID")
	c.SetParamNames("placeID")
	c.SetParamValues("1")

	// Setting up service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setting up Env
	t.Setenv("BASE_URL", "localhost:8080")

	// Setting up input and output
	placeID := 1

	placeDetail := Detail{
		ID:            1,
		Name:          "test_name_place",
		Image:         "test_image_place",
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
	mockService.On("GetDetail", placeID).Return(&placeDetail, nil)

	// Test Fields
	if assert.NoError(t, h.GetDetail(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetDetailWithPlaceIdString(t *testing.T) {
	// Setting up echo router
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:placeID")
	c.SetParamNames("placeID")
	c.SetParamValues("satu")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Expectation
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: "input validation error",
		Errors: []string{
			"placeID must be number",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	response := h.GetDetail(c)
	util.ErrorHandler(response, c)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestService_GetPlaceListWithPlaceIdBelowOne(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:placeID")
	c.SetParamNames("placeID")
	c.SetParamValues("0")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input
	placeID := 0

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"placeID must be above 0"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}
	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var placeDetail Detail
	mockService.On("GetDetail", placeID).Return(&placeDetail, errorFromService)

	response := h.GetDetail(c)
	util.ErrorHandler(response, c)

	// Test
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestService_GetPlaceListWithInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/place/:placeID")
	c.SetParamNames("placeID")
	c.SetParamValues("10")

	// Setup service
	mockService := new(MockService)
	h := NewHandler(mockService)

	// Define input and output
	placeID := 10

	errorFromService := errors.Wrap(ErrInternalServerError, "test error")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	// Excpectation
	var placeDetail Detail
	mockService.On("GetDetail", placeID).Return(&placeDetail, errorFromService)

	response := h.GetDetail(c)
	util.ErrorHandler(response, c)

	// Tes
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
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

	response := h.GetPlacesListWithPagination(c)
	util.ErrorHandler(response, c)

	// Tes
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

	response := h.GetPlacesListWithPagination(c)
	util.ErrorHandler(response, c)

	// Tes
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
	response := h.GetPlacesListWithPagination(c)
	util.ErrorHandler(response, c)

	// Tes
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

func TestHandler_GetListReviewAndRatingWithPaginationSuccess(t *testing.T) {
	e := echo.New()

	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	q.Set("latest", "true")
	q.Set("rating", "true")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/place/1/review",
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
	}

	t.Setenv("BASE_URL", "localhost:8080")

	listReview := ListReview{
		Reviews: []Review{
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
		TotalCount: 10,
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

	mockService.On("GetListReviewAndRatingWithPagination", params).Return(&listReview, pagination, nil)

	if assert.NoError(t, h.GetListReviewAndRatingWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListReviewAndRatingWithPaginationPlaceIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	q.Set("latest", "true")
	q.Set("rating", "true")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
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
	util.ErrorHandler(h.GetListReviewAndRatingWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListReviewAndRatingWithPaginationQueryParamsError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "asd")
	q.Set("page", "asd")
	q.Set("latest", "asd")
	q.Set("rating", "asd")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
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
			"latest parameter should be boolean type",
			"rating parameter should be boolean type",
			"limit should be positive integer",
			"page should be positive integer",
		},
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	util.ErrorHandler(h.GetListReviewAndRatingWithPagination(ctx), ctx)

	// Tes
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListReviewAndRatingWithPaginationLimitError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("limit", "110")
	q.Set("page", "1")
	q.Set("latest", "true")
	q.Set("rating", "true")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListReviewRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/place/1/review",
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
	}

	errorFromService := errors.Wrap(ErrInputValidationError, strings.Join([]string{"limit should be 1 - 100"}, ","))
	errList, errMessage := util.ErrorUnwrap(errorFromService)
	expectedResponse := util.APIResponse{
		Status:  http.StatusBadRequest,
		Message: errMessage,
		Errors:  errList,
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listReview ListReview
	var pagination util.Pagination
	mockService.On("GetListReviewAndRatingWithPagination", params).Return(&listReview, pagination, errorFromService)
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
	q.Set("limit", "110")
	q.Set("page", "1")
	q.Set("latest", "true")
	q.Set("rating", "true")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	params := ListReviewRequest{
		Limit:   110,
		Page:    1,
		Path:    "/api/v1/place/1/review",
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
	}

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJSON, _ := json.Marshal(expectedResponse)

	var listReview ListReview
	var pagination util.Pagination
	mockService.On("GetListReviewAndRatingWithPagination", params).Return(&listReview, pagination, internalServerError)
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
	q.Set("latest", "")
	q.Set("rating", "")
	req := httptest.NewRequest(http.MethodGet, "/review?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/review")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)

	paramsDefault := ListReviewRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/v1/place/1/review",
		PlaceID: 1,
		Latest:  true,
		Rating:  false,
	}

	t.Setenv("BASE_URL", "localhost:8080")

	listReview := ListReview{
		Reviews: []Review{
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
		TotalCount: 10,
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

	mockService.On("GetListReviewAndRatingWithPagination", paramsDefault).Return(&listReview, pagination, nil)

	if assert.NoError(t, h.GetListReviewAndRatingWithPagination(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJSON), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}
