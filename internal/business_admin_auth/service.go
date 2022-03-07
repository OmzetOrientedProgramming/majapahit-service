package businessadminauth

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type Service interface {
	RegisterBusinessAdmin() (bool, error)
}

type service struct {
	repo Repo
}

func (s service) RegisterBusinessAdmin() (bool, error) {
	panic("implement me")
}
