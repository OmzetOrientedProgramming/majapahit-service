package booking

// Service interface consisted function can be used by service
type Service interface {
	GetDetail(bookingID int) (*Detail, error)
}

type service struct {
	repo Repo
}

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo}
}

func (s *service) GetDetail(bookingID int) (*Detail, error) {
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
		panic(err.Error())
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
