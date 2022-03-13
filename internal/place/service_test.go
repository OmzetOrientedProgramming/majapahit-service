package place

import (
	"testing"

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
