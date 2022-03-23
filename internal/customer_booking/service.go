package customerbooking

import (
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
	listCustomerBooking, err := s.repo.GetListCustomerBookingWithPagination(params)
	if err != nil {
		return nil, nil, err
	}
	pagination := util.GeneratePagination(listCustomerBooking.TotalCount, params.Limit, params.Page, params.Path)

	return listCustomerBooking, &pagination, err
}
