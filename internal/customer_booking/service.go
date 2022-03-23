package customerbooking

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	GetListCustomerBookingWithPagination(params ListRequest) (*List, *util.Pagination, error)
}

type service struct {
	repo Repo
}

func (s service) GetListCustomerBookingWithPagination(params ListRequest) (*List, *util.Pagination, error) {
	var errorList []string;

	if params.State == 0{
		params.State = 1
	}

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
	
	listCustomerBooking, err := s.repo.GetListCustomerBookingWithPagination(params)
	if err != nil {
		return nil, nil, err
	}
	pagination := util.GeneratePagination(listCustomerBooking.TotalCount, params.Limit, params.Page, params.Path)

	return listCustomerBooking, &pagination, err
}
