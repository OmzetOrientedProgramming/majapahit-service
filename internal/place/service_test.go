package place

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetDetail(placeID int) (*Detail, error) {
	args := m.Called(placeID)
	ret := args.Get(0).(Detail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetAverageRatingAndReviews(placeID int) (*AverageRatingAndReviews, error) {
	args := m.Called(placeID)
	ret := args.Get(0).(AverageRatingAndReviews)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, error) {
	args := m.Called(params)
	ret := args.Get(0).(ListReview)
	return &ret, args.Error(1)
}

func TestService_GetDetailSuccess(t *testing.T) {
	placeID := 1
	placeDetail := Detail{
		ID:          1,
		Name:        "test_name_place",
		Image:       "test_image_place",
		Address:     "test_address_place",
		Description: "test_description_place",
		OpenHour:    "08:00",
		CloseHour:   "16:00",
	}

	averageRatingAndReviews := AverageRatingAndReviews{
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

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", placeID).Return(placeDetail, nil)
	mockRepo.On("GetAverageRatingAndReviews", placeID).Return(averageRatingAndReviews, nil)

	placeDetailResult, err := mockService.GetDetail(placeID)
	mockRepo.AssertExpectations(t)

	placeDetail.AverageRating = averageRatingAndReviews.AverageRating
	placeDetail.ReviewCount = averageRatingAndReviews.ReviewCount

	placeDetail.Reviews = make([]UserReview, 2)
	placeDetail.Reviews[0].User = averageRatingAndReviews.Reviews[0].User
	placeDetail.Reviews[0].Rating = averageRatingAndReviews.Reviews[0].Rating
	placeDetail.Reviews[0].Content = averageRatingAndReviews.Reviews[0].Content

	placeDetail.Reviews[1].User = averageRatingAndReviews.Reviews[1].User
	placeDetail.Reviews[1].Rating = averageRatingAndReviews.Reviews[1].Rating
	placeDetail.Reviews[1].Content = averageRatingAndReviews.Reviews[1].Content

	assert.Equal(t, &placeDetail, placeDetailResult)
	assert.NotNil(t, placeDetailResult)
	assert.NoError(t, err)
}

func TestService_GetDetailWrongInput(t *testing.T) {
	// Define input
	placeID := -1

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	placeDetail, err := mockService.GetDetail(placeID)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, placeDetail)
}

func TestService_GetDetailFailedCalledGetDetail(t *testing.T) {
	placeID := 1
	var placeDetail Detail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", placeID).Return(placeDetail, ErrInternalServerError)

	placeDetailResult, err := mockService.GetDetail(placeID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailResult)
}

func TestService_GetDetailFailedCalledGetAverageRatingAndReviews(t *testing.T) {
	placeID := 1
	placeDetail := Detail{}
	averageRatingAndReviews := AverageRatingAndReviews{}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetDetail", placeID).Return(placeDetail, nil)
	mockRepo.On("GetAverageRatingAndReviews", placeID).Return(averageRatingAndReviews, ErrInternalServerError)

	placeDetailResult, err := mockService.GetDetail(placeID)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailResult)
}
func (m *MockRepository) GetPlacesListWithPagination(params PlacesListRequest) (*PlacesList, error) {
	args := m.Called(params)
	ret := args.Get(0).(PlacesList)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetPlaceRatingAndReviewCountByPlaceID(placeID int) (*PlacesRatingAndReviewCount, error) {
	args := m.Called(placeID)
	ret := args.Get(0).(PlacesRatingAndReviewCount)
	return &ret, args.Error(1)
}

func TestService_GetPlaceListWithPaginationSuccess(t *testing.T) {
	// Define input and output
	placeList := PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test name",
				Description: "test description",
				Address:     "test address",
			},
			{
				ID:          2,
				Name:        "test name 2",
				Description: "test description 2",
				Address:     "test address 2",
			},
		},
		TotalCount: 2,
	}

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	ratingAndReview := []PlacesRatingAndReviewCount{
		{
			Rating:      5.0,
			ReviewCount: 10,
		},
		{
			Rating:      5.0,
			ReviewCount: 20,
		},
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetPlacesListWithPagination", params).Return(placeList, nil)
	mockRepo.On("GetPlaceRatingAndReviewCountByPlaceID", placeList.Places[0].ID).Return(ratingAndReview[0], nil)
	mockRepo.On("GetPlaceRatingAndReviewCountByPlaceID", placeList.Places[1].ID).Return(ratingAndReview[0], nil)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)
	mockRepo.AssertExpectations(t)

	placeList.Places[0].Rating = ratingAndReview[0].Rating
	placeList.Places[0].ReviewCount = ratingAndReview[0].ReviewCount

	placeList.Places[1].Rating = ratingAndReview[1].Rating
	placeList.Places[1].ReviewCount = ratingAndReview[1].ReviewCount

	assert.Equal(t, &placeList, placeListResult)
	assert.NotNil(t, placeListResult)
	assert.NoError(t, err)
}

func TestService_GetPlaceListWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	placeList := PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test name",
				Description: "test description",
				Address:     "test address",
			},
			{
				ID:          2,
				Name:        "test name 2",
				Description: "test description 2",
				Address:     "test address 2",
			},
		},
		TotalCount: 2,
	}

	params := PlacesListRequest{
		Limit: 0,
		Page:  0,
		Path:  "/api/testing",
	}

	ratingAndReview := []PlacesRatingAndReviewCount{
		{
			Rating:      5.0,
			ReviewCount: 10,
		},
		{
			Rating:      5.0,
			ReviewCount: 20,
		},
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	paramsDefault := PlacesListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Expectation
	mockRepo.On("GetPlacesListWithPagination", paramsDefault).Return(placeList, nil)
	mockRepo.On("GetPlaceRatingAndReviewCountByPlaceID", placeList.Places[0].ID).Return(ratingAndReview[0], nil)
	mockRepo.On("GetPlaceRatingAndReviewCountByPlaceID", placeList.Places[1].ID).Return(ratingAndReview[0], nil)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)
	mockRepo.AssertExpectations(t)

	placeList.Places[0].Rating = ratingAndReview[0].Rating
	placeList.Places[0].ReviewCount = ratingAndReview[0].ReviewCount

	placeList.Places[1].Rating = ratingAndReview[1].Rating
	placeList.Places[1].ReviewCount = ratingAndReview[1].ReviewCount

	assert.Equal(t, &placeList, placeListResult)
	assert.NotNil(t, placeListResult)
	assert.NoError(t, err)
}

func TestService_GetPlaceListWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	params := PlacesListRequest{
		Limit: 101,
		Page:  0,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, placeListResult)
}

func TestService_GetPlaceListWithPaginationFailedCalledGetPlacesListWithPagination(t *testing.T) {
	// Define input and output
	var placesList PlacesList

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetPlacesListWithPagination", params).Return(placesList, ErrInternalServerError)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeListResult)
}

func TestService_GetPlaceListWithPaginationFailedCalledGetPlaceRatingAndReviewCountByPlaceID(t *testing.T) {
	// Define input and output
	placeList := PlacesList{
		Places: []Place{
			{
				ID:          1,
				Name:        "test name",
				Description: "test description",
				Address:     "test address",
			},
			{
				ID:          2,
				Name:        "test name 2",
				Description: "test description 2",
				Address:     "test address 2",
			},
		},
		TotalCount: 2,
	}

	var ratingAndReview PlacesRatingAndReviewCount

	params := PlacesListRequest{
		Limit: 10,
		Page:  1,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetPlacesListWithPagination", params).Return(placeList, nil)
	mockRepo.On("GetPlaceRatingAndReviewCountByPlaceID", placeList.Places[0].ID).Return(ratingAndReview, ErrInternalServerError)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeListResult)
}

func TestService_GetPlaceListWithPaginationFailedURLIsEmpty(t *testing.T) {
	// Define input
	params := PlacesListRequest{
		Limit: 100,
		Page:  0,
		Path:  "",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	placeListResult, _, err := mockService.GetPlaceListWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, placeListResult)
}

func TestService_GetListReviewAndRatingWithPaginationSuccess(t *testing.T) {
	// Define input and output
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

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetListReviewAndRatingWithPagination", params).Return(listReview, nil)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listReview, listReviewResult)
	assert.NotNil(t, listReviewResult)
	assert.NoError(t, err)
}

func TestService_GetListReviewAndRatingWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
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

	params := ListReviewRequest{
		Limit:   0,
		Page:    0,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	paramsDefault := ListReviewRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Expectation
	mockRepo.On("GetListReviewAndRatingWithPagination", paramsDefault).Return(listReview, nil)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listReview, listReviewResult)
	assert.NotNil(t, listReviewResult)
	assert.NoError(t, err)
}

func TestService_GetListReviewAndRatingWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	params := ListReviewRequest{
		Limit:   101,
		Page:    1,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listReviewResult)
}

func TestService_GetListReviewAndRatingWithPaginationFailedURLEmpty(t *testing.T) {
	// Define input
	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listReviewResult)
}

func TestService_GetListReviewAndRatingWithPaginationFailedPlaceIDNegative(t *testing.T) {
	// Define input
	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		PlaceID: -1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listReviewResult)
}

func TestService_GetListReviewAndRatingWithPaginationFailedCalledRepo(t *testing.T) {
	// Define input and output
	var listReview ListReview

	params := ListReviewRequest{
		Limit:   10,
		Page:    1,
		PlaceID: 1,
		Latest:  true,
		Rating:  true,
		Path:    "/api/testing",
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetListReviewAndRatingWithPagination", params).Return(listReview, ErrInternalServerError)

	// Test
	listReviewResult, _, err := mockService.GetListReviewAndRatingWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listReviewResult)
}
