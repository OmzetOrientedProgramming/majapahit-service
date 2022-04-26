package cloudinary

import (
	"context"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Repo contains all the function that available of this repo package
type Repo interface {
	UploadFile(fileContent, folderName, fileName string) (string, error)
}

type repo struct {
	cloudinary *cloudinary.Cloudinary
}

// NewRepo for initialize repo
func NewRepo(cloudName, apiKey, apiSecret string) Repo {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		logrus.Fatalf("failed to initialize cloudinary: %v", err)
	}
	return &repo{
		cloudinary: cld,
	}
}

func (r repo) UploadFile(fileContent, folderName, fileName string) (string, error) {
	resp, err := r.cloudinary.Upload.Upload(context.Background(), fileContent, uploader.UploadParams{
		Folder:       folderName,
		PublicID:     fileName,
		ResourceType: "auto",
	})
	if err != nil {
		return "", errors.Wrap(ErrInternalServer, err.Error())
	}
	if resp.Error.Message != "" {
		return "", errors.Wrap(ErrInternalServer, resp.Error.Message)
	}
	return resp.SecureURL, nil
}
