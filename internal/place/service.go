package place

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service will contain all the function that can be used by service
type Service interface {
	GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error)
	GetPlaceDetail(placeId int) (*PlaceDetail, error)
}

type service struct {
	repo Repo
}

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo}
}

func (s *service) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	placeDetail, err := s.repo.GetPlaceDetail(placeId)
	if err != nil {
		return nil, err
	}

	averageRatingAndReviews, err := s.repo.GetAverageRatingAndReviews(placeId)
	if err != nil {
		return nil, err
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

func (s service) GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error) {
	var errorList []string

	if params.Page == 0 {
		params.Page = util.DefaultPage
	}

	if params.Limit == 0 {
		params.Limit = util.DefaultLimit
	}

	if params.Limit > util.MaxLimit {
		errorList = append(errorList, "limit should be 1 - 100")
	}

	if params.Path == "" {
		errorList = append(errorList, "path is required for pagination")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	placeList, err := s.repo.GetPlacesListWithPagination(params)
	if err != nil {
		return nil, nil, err
	}

	for i := range placeList.Places {
		ratingAndReviewCountRetrieved, errRating := s.repo.GetPlaceRatingAndReviewCountByPlaceID(placeList.Places[i].ID)
		if errRating != nil {
			return nil, nil, errRating
		}

		placeList.Places[i].Rating = ratingAndReviewCountRetrieved.Rating
		placeList.Places[i].ReviewCount = ratingAndReviewCountRetrieved.ReviewCount
	}

	pagination := util.GeneratePagination(placeList.TotalCount, params.Limit, params.Page, params.Path)

	return placeList, &pagination, err
}
