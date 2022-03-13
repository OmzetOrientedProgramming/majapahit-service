package checkup

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	GetApplicationCheckUp() (bool, error)
}

type service struct {
	repo Repo
}

func (s service) GetApplicationCheckUp() (bool, error) {
	isUp, err := s.repo.GetApplicationCheckUp()
	if err != nil {
		return false, err
	}

	return isUp, nil
}
