package upload

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/cloudinary"
)

// NewService for initialize service
func NewService(cloudinaryRepo cloudinary.Repo) Service {
	return &service{
		cloudinaryRepo: cloudinaryRepo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	UploadProfilePicture(params FileRequest) (*FileResponse, error)
}

type service struct {
	cloudinaryRepo cloudinary.Repo
}

func (s service) UploadProfilePicture(params FileRequest) (*FileResponse, error) {
	var errorList []string

  if params.File == "" {
    errorList = append(errorList, "File diperlukan")
  }

  if params.CustomerName == "" {
    errorList = append(errorList, "CustomerName diperlukan")
  }

  if len(errorList) > 0 {
    return nil, errors.Wrap(ErrInputValidation, strings.Join(errorList, ";"))
  }

  url, err := s.cloudinaryRepo.UploadFile(params.File, "Profile Picture", fmt.Sprintf("%s-Profile-Picture", params.CustomerName))
  if err != nil {
    return nil, err
  }

  response := FileResponse{
    URL: url,
  }

  return &response, nil
}
