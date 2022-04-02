package booking

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Service will contain all the function that can be used by service
type Service interface {
	GetDetail(bookingID int) (*Detail, error)
	UpdateBookingStatus(bookingID int, newStatus int) error
	GetMyBookingsOngoing(localID string) (*[]Booking, error)
	GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error)
}

type service struct {
	repo Repo
}

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetDetail(bookingID int) (*Detail, error) {
	errorList := []string{}

	if bookingID <= 0 {
		errorList = append(errorList, "bookingID must be above 0")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	bookingDetail, err := s.repo.GetDetail(bookingID)
	if err != nil {
		return nil, err
	}

	ticketPriceWrapper, err := s.repo.GetTicketPriceWrapper(bookingID)
	if err != nil {
		return nil, err
	}

	itemsWrapper, err := s.repo.GetItemWrapper(bookingID)
	if err != nil {
		return nil, err
	}

	totalPriceTicket := ticketPriceWrapper.Price * float64(bookingDetail.Capacity)
	totalPrice := totalPriceTicket + bookingDetail.TotalPriceItem

	bookingDetail.TotalPriceTicket = totalPriceTicket
	bookingDetail.TotalPrice = totalPrice

	bookingDetail.Items = make([]ItemDetail, 0)
	for _, item := range itemsWrapper.Items {
		bookingDetail.Items = append(bookingDetail.Items, item)
	}

	return bookingDetail, nil
}

func (s *service) UpdateBookingStatus(bookingID int, newStatus int) error {
	errorList := []string{}

	if bookingID <= 0 {
		errorList = append(errorList, "bookingID must be above 0")
	}

	if newStatus < 0 {
		errorMessage := fmt.Sprintf("there are no status: %d", newStatus)
		errorList = append(errorList, errorMessage)
	}

	if len(errorList) > 0 {
		return errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	err := s.repo.UpdateBookingStatus(bookingID, newStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetMyBookingsOngoing(localID string) (*[]Booking, error) {
	errorList := []string{}

	if localID == "" {
		errorList = append(errorList, "localID cannot be empty")
	}

	if len(errorList) > 0 {
		return nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	myBookingsOngoing, err := s.repo.GetMyBookingsOngoing(localID)
	if err != nil {
		return nil, err
	}

	return myBookingsOngoing, nil
}

func (s service) GetMyBookingsPreviousWithPagination(localID string, params BookingsListRequest) (*List, *util.Pagination, error) {
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

	myBookingsPrevious, err := s.repo.GetMyBookingsPreviousWithPagination(localID, params)
	if err != nil {
		return nil, nil, err
	}

	pagination := util.GeneratePagination(myBookingsPrevious.TotalCount, params.Limit, params.Page, params.Path)

	return myBookingsPrevious, &pagination, err
}
