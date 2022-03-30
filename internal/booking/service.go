package booking

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
	panic("Not yet implemented!")
}
