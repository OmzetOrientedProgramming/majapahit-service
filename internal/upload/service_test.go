package upload

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CloudinaryMockRepository struct {
	mock.Mock
}

func (c *CloudinaryMockRepository) UploadFile(fileContent, folderName, fileName string) (string, error){
	args := c.Called(fileContent, folderName, fileName)

	return args.String(0), args.Error(1)
}

func TestService_UploadProfilePicture(t * testing.T) {
	mockRepo := new(CloudinaryMockRepository)
	mockService := NewService(mockRepo)

	t.Run("Upload profile picture done successfully", func(t *testing.T){
		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		customerName := "Testing User"

		params := FileRequest {
			File: file,
			CustomerName: customerName,
		}		

		uplaodedURL := "https://res.cloudinary.com/wave-ppl/image/upload/v1652093802/Profile%20Picture/Mario%20Serano-Profile-Picture.png"
		expectedResponse := FileResponse {
			URL: uplaodedURL,
		}

		mockRepo.On("UploadFile", params.File, "Profile Picture", fmt.Sprintf("%s-Profile-Picture", params.CustomerName)).Return(uplaodedURL, nil)
		response, err := mockService.UploadProfilePicture(params)

		assert.Equal(t, response, &expectedResponse)
		assert.NotNil(t, response)
		assert.NoError(t, err)
	})

	t.Run("File is empty", func(t *testing.T){
		file := ""
		customerName := "Testing User"

		params := FileRequest {
			File: file,
			CustomerName: customerName,
		}		

		response, err := mockService.UploadProfilePicture(params)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
		assert.Nil(t, response)
	})

	t.Run("Customer name is empty", func(t *testing.T){
		file := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="
		customerName := ""

		params := FileRequest {
			File: file,
			CustomerName: customerName,
		}		

		response, err := mockService.UploadProfilePicture(params)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
		assert.Nil(t, response)
	})

	t.Run("Cloudinary error", func(t *testing.T){
		file := "fix error"
		customerName := "Testing User"

		params := FileRequest {
			File: file,
			CustomerName: customerName,
		}		

		expectedError := errors.Wrap(ErrInputValidation, "error gess")

		mockRepo.On("UploadFile", params.File, "Profile Picture", fmt.Sprintf("%s-Profile-Picture", params.CustomerName)).Return("", expectedError)
		response, err := mockService.UploadProfilePicture(params)

		assert.Equal(t, err, expectedError)
		assert.Nil(t, response)
	})
}

