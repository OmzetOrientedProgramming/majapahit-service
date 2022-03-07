package item

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func (m *MockService) GetListItem(placeID int, name string) (*ListItem, error) {
	args := m.Called(placeID, name)
	listItem := args.Get(0).(*ListItem)
	return listItem, args.Error(1)
}

func (m *MockService) GetItemByID(placeID int, itemID int ) (*Item, error) {
	args := m.Called(placeID, itemID)
	item := args.Get(0).(*Item)
	return item, args.Error(1)
}

func TestHandler_GetListItem(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	placeID := 1
	name := ""

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
	}

	expectedResponse := util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":     listItem.Items,
		},
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	mockService.On("GetListItem", placeID, name).Return(&listItem, nil)

	// Tes
	if assert.NoError(t, h.GetListItem(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
	}
}

func TestHandler_GetListItemPlaceIDError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
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
	assert.NoError(t, h.GetListItem(ctx))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestHandler_GetListItemInternalServerError(t *testing.T) {
	// Setup echo
	e := echo.New()

	// import "net/url"
	q := make(url.Values)
	q.Set("name", "")
	req := httptest.NewRequest(http.MethodGet, "/catalog?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/:placeID/catalog")
	ctx.SetParamNames("placeID")
	ctx.SetParamValues("1")

	mockService := new(MockService)
	h := NewHandler(mockService)
	placeID := 1
	name := ""

	// Setup Env
	t.Setenv("BASE_URL", "localhost:8080")

	internalServerError := errors.Wrap(ErrInternalServerError, "test")
	expectedResponse := util.APIResponse{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
	}

	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// Excpectation
	var listItem ListItem
	mockService.On("GetListItem", placeID, name).Return(&listItem, internalServerError)
	
	// Tes
	assert.NoError(t, h.GetListItem(ctx))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, string(expectedResponseJson), strings.TrimSuffix(rec.Body.String(), "\n"))
}