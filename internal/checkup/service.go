package checkup

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

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
