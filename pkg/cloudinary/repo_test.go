package cloudinary

import (
	"errors"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRepo_UploadFile(t *testing.T) {
	tests := map[string]struct {
		expectedFile string
		wantError    error
	}{
		"success": {
			expectedFile: "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg==",
			wantError:    nil,
		},
		"unsopperted file format": {
			expectedFile: "1",
			wantError:    ErrInternalServer,
		},
		"upload failure": {
			expectedFile: "data:",
			wantError:    ErrInternalServer,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_ = godotenv.Load("../../.env")

			expectedFolderName := "mockFolderName"
			expectedFileName := "mockFileName"

			resp, err := NewRepo(
				os.Getenv("CLOUDINARY_CLOUD_NAME"),
				os.Getenv("CLOUDINARY_API_KEY"),
				os.Getenv("CLOUDINARY_API_SECRET")).
				UploadFile(test.expectedFile, expectedFolderName, expectedFileName)

			if test.wantError != nil {
				assert.True(t, errors.Is(err, test.wantError))
			} else {
				assert.NotEmpty(t, resp)
				assert.Nil(t, err)
			}
		})
	}
}
