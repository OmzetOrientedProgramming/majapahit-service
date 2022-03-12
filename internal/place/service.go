package place

type Service interface {
	GetPlaceDetail(placeId int) (*PlaceDetail, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetPlaceDetail(placeId int) (*PlaceDetail, error) {
	panic("implement this")
}
