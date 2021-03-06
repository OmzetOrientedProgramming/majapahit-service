package customer

import (
	"strings"
	"time"

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
	PutEditCustomer(body EditCustomerRequest) error
	RetrieveCustomerProfile(userID int) (*Profile, error)
}

type service struct {
	repo Repo
}

func (s service) PutEditCustomer(body EditCustomerRequest) error {
	var errorList []string

	if body.Name == "" {
		errorList = append(errorList, "Name diperlukan")
	}

	if len(body.Name) < util.MinimumNameLength {
		errorList = append(errorList, "Name terlalu pendek")
	}

	if body.ProfilePicture == "" {
		errorList = append(errorList, "Profile picture diperlukan")
	}

	if body.DateOfBirth.IsZero() {
		errorList = append(errorList, "Date of birth diperlukan")
  } else {
    currentTime := time.Now()
    timeDiff := currentTime.Sub(body.DateOfBirth)
    if timeDiff < 0 {
      errorList = append(errorList, "Date of birth tidak boleh melebihi tanggal saat ini")
    }
	}

	if !(body.Gender == util.GenderMale || body.Gender == util.GenderFemale) {
		errorList = append(errorList, "Gender tidak sesuai")
	}

	if len(errorList) > 0 {
		return errors.Wrap(ErrInputValidation, strings.Join(errorList, ";"))
	}

	err := s.repo.PutEditCustomer(body)
	if err != nil {
		return err
	}

	return nil

}

// RetrieveCustomerProfile is called to get customer profile through repository
func (s service) RetrieveCustomerProfile(userID int) (*Profile, error) {

	customerProfile, err := s.repo.RetrieveCustomerProfile(userID)

	if err != nil {
		return nil, err
	}

	return customerProfile, nil
}
