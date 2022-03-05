package item

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type Service interface {
	GetListItem(placeID int, name string) (*ListItem, error)
}

type service struct {
	repo Repo
}

func (s service) GetListItem(placeID int, name string) (*ListItem, error) {
	listItem, err := s.repo.GetListItem(placeID, name)

	return listItem, err
}