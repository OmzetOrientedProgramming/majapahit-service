package place

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service will contain all the function that can be used by service
type Service interface {
	GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error)
	GetDetail(placeID int) (*Detail, error)
	GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, *util.Pagination, error)
}

type service struct {
	repo Repo
}

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo}
}

func (s *service) GetDetail(placeID int) (*Detail, error) {
	errorList := []string{}

	if placeID <= 0 {
		errorList = append(errorList, "placeID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ";"))
	}

	placeDetail, err := s.repo.GetDetail(placeID)
	if err != nil {
		return nil, err
	}

	averageRatingAndReviews, err := s.repo.GetAverageRatingAndReviews(placeID)
	if err != nil {
		return nil, err
	}

	placeDetail.AverageRating = averageRatingAndReviews.AverageRating
	placeDetail.ReviewCount = averageRatingAndReviews.ReviewCount

	placeDetail.Reviews = make([]UserReview, 0)
	for _, i := range averageRatingAndReviews.Reviews {
		placeDetail.Reviews = append(placeDetail.Reviews, i)
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

	switch params.Sort {
	case "", "recommended":
		params.Sort = "recommended"
	case "distance", "popularity":
	default:
		errorList = append(errorList, "invalid sort value")
	}

	for _, people := range params.People {
		if people != "1" && people != "2-4" && people != "5-10" && people != "10" {
			errorList = append(errorList, "invalid people filter value")
		}
	}

	for _, price := range params.Price {
		if price != "16000" && price != "16000-40000" && price != "40000-100000" && price != "100000" {
			errorList = append(errorList, "invalid price filter value")
		}
	}

	for _, rating := range params.Rating {
		if rating < 1 || rating > 5 {
			errorList = append(errorList, "invalid rating filter value")
		}
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ";"))
	}

	placeList, err := s.repo.GetPlacesListWithPagination(params)
	if err != nil {
		return nil, nil, err
	}

	//for i := range placeList.Places {
	//	ratingAndReviewCountRetrieved, errRating := s.repo.GetPlaceRatingAndReviewCountByPlaceID(placeList.Places[i].ID)
	//	if errRating != nil {
	//		return nil, nil, errRating
	//	}
	//
	//	placeList.Places[i].Rating = ratingAndReviewCountRetrieved.Rating
	//	placeList.Places[i].ReviewCount = ratingAndReviewCountRetrieved.ReviewCount
	//}

	pagination := util.GeneratePagination(placeList.TotalCount, params.Limit, params.Page, params.Path)

	return placeList, &pagination, err
}

func (s service) GetListReviewAndRatingWithPagination(params ListReviewRequest) (*ListReview, *util.Pagination, error) {
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

	if params.PlaceID <= 0 {
		errorList = append(errorList, "placeID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ";"))
	}

	listReview, err := s.repo.GetListReviewAndRatingWithPagination(params)
	if err != nil {
		return nil, nil, err
	}

	pagination := util.GeneratePagination(listReview.TotalCount, params.Limit, params.Page, params.Path)

	return listReview, &pagination, err
}
