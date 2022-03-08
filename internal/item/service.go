package item

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

type Service interface {
	GetListItem(placeID int, name string) (*ListItem, error)
	GetItemByID(placeID int, itemID int) (*Item, error)
}

type service struct {
	repo Repo
}

func (s service) GetListItem(placeID int, name string) (*ListItem, error) {
	listItem, err := s.repo.GetListItem(placeID, name)

	if err != nil {
		return nil, err
	}

	return listItem, err
}

func (s service) GetItemByID(placeID int, itemID int) (*Item, error) {
	item, err := s.repo.GetItemByID(placeID, itemID)

	if err != nil {
		return nil, err
	}

	return item, err
}