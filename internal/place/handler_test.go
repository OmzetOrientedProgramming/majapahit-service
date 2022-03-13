package place

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
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
