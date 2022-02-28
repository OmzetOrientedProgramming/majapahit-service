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
	var item, error = s.repo.GetItem()
	// Do something
	return item, error
}