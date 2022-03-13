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

func (m *MockRepository) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	args := m.Called(placeId)
	ret := args.Get(0).(PlaceDetail)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetAverageRatingAndReviews(placeId int) (*AverageRatingAndReviews, error) {
	args := m.Called(placeId)
	ret := args.Get(0).(AverageRatingAndReviews)
	return &ret, args.Error(1)
}

func TestService_GetPlaceDetailSuccess(t *testing.T) {
	placeId := 1
	placeDetail := PlaceDetail{
		ID:          1,
		Name:        "test_name_place",
		Image:       "test_image_place",
		Distance:    200,
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

	mockRepo.On("GetPlaceDetail", placeId).Return(placeDetail, nil)
	mockRepo.On("GetAverageRatingAndReviews", placeId).Return(averageRatingAndReviews, nil)

	placeDetailResult, err := mockService.GetPlaceDetail(placeId)
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

func TestService_GetPlaceDetailWrongInput(t *testing.T) {
	// Define input
	placeId := -1

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Test
	placeDetail, err := mockService.GetPlaceDetail(placeId)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, placeDetail)
}

func TestService_GetPlaceDetailFailedCalledGetPlaceDetail(t *testing.T) {
	placeId := 1
	var placeDetail PlaceDetail

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetPlaceDetail", placeId).Return(placeDetail, ErrInternalServerError)

	placeDetailResult, err := mockService.GetPlaceDetail(placeId)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, placeDetailResult)
}

func TestService_GetPlaceDetailFailedCalledGetAverageRatingAndReviews(t *testing.T) {
	placeId := 1
	placeDetail := PlaceDetail{}
	averageRatingAndReviews := AverageRatingAndReviews{}

	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetPlaceDetail", placeId).Return(placeDetail, nil)
	mockRepo.On("GetAverageRatingAndReviews", placeId).Return(averageRatingAndReviews, ErrInternalServerError)

	placeDetailResult, err := mockService.GetPlaceDetail(placeId)
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
