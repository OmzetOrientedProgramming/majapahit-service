package item

import (
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	GetListItemWithPagination(params ListItemRequest) (*ListItem, *util.Pagination, error)
	GetItemByID(placeID int, itemID int) (*Item, error)
}

type service struct {
	repo Repo
}

func (s service) GetListItemWithPagination(params ListItemRequest) (*ListItem, *util.Pagination, error) {
	var errorList []string

	if strings.Contains(params.Name, "+") {
		params.Name = strings.ReplaceAll(params.Name, "+", " ")
	}

	if params.Page == 0 {
		params.Page = util.DefaultPage
	}

	if params.Limit == 0 {
		params.Limit = util.DefaultLimit
	}

	if params.Limit > util.MaxLimit {
		errorList = append(errorList, "limit should be 1 - 100")
	}

	if params.Path == "" {
		errorList = append(errorList, "path is required for pagination")
	}

	if len(errorList) > 0 {
		return nil, nil, errors.Wrap(ErrInputValidationError, strings.Join(errorList, ","))
	}

	listItem, err := s.repo.GetListItemWithPagination(params)

	if err != nil {
		return nil, nil, err
	}

	pagination := util.GeneratePagination(listItem.TotalCount, params.Limit, params.Page, params.Path)
	return listItem, &pagination, err
}

func (s service) GetItemByID(placeID int, itemID int) (*Item, error) {
	item, err := s.repo.GetItemByID(placeID, itemID)

	if err != nil {
		return nil, err
	}

	return item, err
}
