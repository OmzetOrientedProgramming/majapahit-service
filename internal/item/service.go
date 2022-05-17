package item

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/cloudinary"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// NewService for initialize service
func NewService(repo Repo, cloudinary cloudinary.Repo) Service {
	return &service{
		repo:       repo,
		cloudinary: cloudinary,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	GetListItemWithPagination(params ListItemRequest) (*ListItem, *util.Pagination, error)
	GetItemByID(placeID int, itemID int) (*Item, error)
	DeleteItemAdminByID(itemID int) error
	UpdateItem(ID int, item Item) error
	CreateItem(userID int, item Item) error
}

type service struct {
	repo       Repo
	cloudinary cloudinary.Repo
}

func (s service) GetListItemWithPagination(params ListItemRequest) (*ListItem, *util.Pagination, error) {
	var errorList []string
	var listItem *ListItem
	var err error

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

	if params.PlaceID != 0 {
		listItem, err = s.repo.GetListItemWithPagination(params)
	} else {
		listItem, err = s.repo.GetListItemAdminWithPagination(params)
	}

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

func (s service) DeleteItemAdminByID(itemID int) error {
	err := s.repo.DeleteItemAdminByID(itemID)

	if err != nil {
		return err
	}

	return nil
}

func (s service) UpdateItem(ID int, item Item) error {
	if !strings.HasPrefix(item.Image, "data:") {
		return fmt.Errorf("string is not a data URI: %w", ErrInputValidationError)
	}
	withoutData := strings.TrimPrefix(item.Image, "data:")

	if !strings.HasPrefix(withoutData, "image/") {
		return fmt.Errorf("string is not an image data URI: %w", ErrInputValidationError)
	}

	imageString := strings.Split(withoutData, ",")[1]

	_, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		logrus.Errorf("image string is not base64: %v", err)
		return fmt.Errorf("image string is not base64: %w", ErrInputValidationError)
	}

	imageURL, err := s.cloudinary.UploadFile(item.Image, "Item Image", fmt.Sprintf("%d-%s", item.ID, item.Name))
	if err != nil {
		return err
	}

	item.Image = imageURL
	if err := s.repo.UpdateItem(ID, item); err != nil {
		return err
	}
	return nil
}

func (s service) CreateItem(userID int, item Item) error {
	if !strings.HasPrefix(item.Image, "data:") {
		return fmt.Errorf("string is not a data URI: %w", ErrInputValidationError)
	}
	withoutData := strings.TrimPrefix(item.Image, "data:")

	if !strings.HasPrefix(withoutData, "image/") {
		return fmt.Errorf("string is not an image data URI: %w", ErrInputValidationError)
	}

	imageString := strings.Split(withoutData, ",")[1]

	_, err := base64.StdEncoding.DecodeString(imageString)
	if err != nil {
		logrus.Errorf("image string is not base64: %v", err)
		return fmt.Errorf("image string is not base64: %w", ErrInputValidationError)
	}

	imageURL, err := s.cloudinary.UploadFile(item.Image, "Item Image", fmt.Sprintf("%s", item.Name))
	if err != nil {
		return err
	}

	item.Image = imageURL
	if err := s.repo.CreateItem(userID, item); err != nil {
		return err
	}

	return nil
}
