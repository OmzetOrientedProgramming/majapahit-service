package place

type Service interface {
	GetPlaceDetail(placeId int) (*PlaceDetail, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	placeDetail, err := s.repo.GetPlaceDetail(placeId)
	if err != nil {

	}

	averageRatingAndReviews, err := s.repo.GetAverageRatingAndReviews(placeId)
	if err != nil {

	}

	placeDetail.AverageRating = averageRatingAndReviews.AverageRating
	placeDetail.ReviewCount = averageRatingAndReviews.ReviewCount

	placeDetail.Reviews = make([]UserReview, 2)
	for i := range averageRatingAndReviews.Reviews {
		placeDetail.Reviews[i].User = averageRatingAndReviews.Reviews[i].User
		placeDetail.Reviews[i].Rating = averageRatingAndReviews.Reviews[i].Rating
		placeDetail.Reviews[i].Content = averageRatingAndReviews.Reviews[i].Content
	}

	return placeDetail, nil
}
