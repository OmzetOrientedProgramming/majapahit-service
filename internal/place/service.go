package place

import (
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"strings"
)

type Service interface {
	GetPlaceListWithPagination(params PlacesListRequest) (*PlacesList, *util.Pagination, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
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
