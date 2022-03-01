package item

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type Service interface {
	GetListItem() (*Item, error)
}

type service struct {
	repo Repo
}

func (s service) GetListItem() (*Item, error) {
	panic("implement me")
}