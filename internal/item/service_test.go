package item

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/cloudinary"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetListItemWithPagination(params ListItemRequest) (*ListItem, error) {
	args := m.Called(params)
	ret := args.Get(0).(ListItem)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetItemByID(placeID int, itemID int) (*Item, error) {
	args := m.Called(placeID, itemID)
	ret := args.Get(0).(Item)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetListItemAdminWithPagination(params ListItemRequest) (*ListItem, error) {
	args := m.Called(params)
	ret := args.Get(0).(ListItem)
	return &ret, args.Error(1)
}

func (m *MockRepository) DeleteItemAdminByID(itemID int) error {
	args := m.Called(itemID)
	return args.Error(0)
}

func (m *MockRepository) UpdateItem(ID int, item Item) error {
	args := m.Called(ID, item)
	return args.Error(0)
}

type MockCloudinary struct {
	isError     bool
	imageString string
}

func (mc MockCloudinary) UploadFile(fileContent, folderName, fileName string) (string, error) {
	if mc.isError {
		return mc.imageString, ErrInternalServerError
	}
	return mc.imageString, nil
}

func TestService_GetListItemByIDWithPaginationSuccess(t *testing.T) {
	// Define input and output
	listItemExpected := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test 1",
				Image:       "test 1",
				Description: "test 1",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test 2",
				Image:       "test 2",
				Description: "test 2",
				Price:       20000,
			},
		},
		TotalCount: 10,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	t.Run("success with place id", func(t *testing.T) {
		params := ListItemRequest{
			Limit:   10,
			Page:    1,
			Path:    "/api/testing",
			PlaceID: 1,
			UserID:  0,
		}
		// Expectation
		mockRepo.On("GetListItemWithPagination", params).Return(listItemExpected, nil)

		// Test
		listItemResult, _, err := mockService.GetListItemWithPagination(params)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, &listItemExpected, listItemResult)
		assert.NotNil(t, listItemResult)
		assert.NoError(t, err)
	})

	t.Run("success with user id", func(t *testing.T) {
		params := ListItemRequest{
			Limit:   10,
			Page:    1,
			Path:    "/api/testing",
			PlaceID: 0,
			UserID:  1,
		}
		// Expectation
		mockRepo.On("GetListItemAdminWithPagination", params).Return(listItemExpected, nil)

		// Test
		listItemResult, _, err := mockService.GetListItemWithPagination(params)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, &listItemExpected, listItemResult)
		assert.NotNil(t, listItemResult)
		assert.NoError(t, err)
	})
}

func TestService_GetListItemByIDWithPaginationSuccessWithParamsName(t *testing.T) {
	// Define input and output
	listItemExpected := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test 1",
				Image:       "test 1",
				Description: "test 1",
				Price:       10000,
			},
		},
		TotalCount: 10,
	}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/testing",
		PlaceID: 1,
		Name:    "test+1",
	}

	newParams := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/testing",
		PlaceID: 1,
		Name:    "test 1",
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	// Expectation
	mockRepo.On("GetListItemWithPagination", newParams).Return(listItemExpected, nil)

	// Test
	listItemResult, _, err := mockService.GetListItemWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestService_GetListItemByIDWithPaginationSuccessWithDefaultParam(t *testing.T) {
	// Define input and output
	listItemExpected := ListItem{
		Items: []Item{
			{
				ID:          1,
				Name:        "test 1",
				Image:       "test 1",
				Description: "test 1",
				Price:       10000,
			},
			{
				ID:          2,
				Name:        "test 2",
				Image:       "test 2",
				Description: "test 2",
				Price:       20000,
			},
		},
		TotalCount: 10,
	}

	params := ListItemRequest{
		Limit:   0,
		Page:    0,
		Path:    "/api/testing",
		PlaceID: 1,
	}

	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	paramsDefault := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/testing",
		PlaceID: 1,
	}
	// Expectation
	mockRepo.On("GetListItemWithPagination", paramsDefault).Return(listItemExpected, nil)

	// Test
	listItemResult, _, err := mockService.GetListItemWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestService_GetListItemWithPaginationFailedLimitExceedMaxLimit(t *testing.T) {
	// Define input
	params := ListItemRequest{
		Limit: 101,
		Page:  0,
		Path:  "/api/testing",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	// Test
	listItemResult, _, err := mockService.GetListItemWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listItemResult)
}

func TestService_GetListItemByIDWithPaginationError(t *testing.T) {
	listItem := ListItem{}

	params := ListItemRequest{
		Limit:   10,
		Page:    1,
		Path:    "/api/testing",
		PlaceID: 1,
	}

	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetListItemWithPagination", params).Return(listItem, ErrInternalServerError)

	// Test
	listItemResult, _, err := mockService.GetListItemWithPagination(params)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listItemResult)
}

func TestService_GetListItemWithPaginationFailedURLIsEmpty(t *testing.T) {
	// Define input
	params := ListItemRequest{
		Limit: 100,
		Page:  0,
		Path:  "",
	}

	// Init mock repo and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	// Test
	listItemResult, _, err := mockService.GetListItemWithPagination(params)

	assert.Equal(t, ErrInputValidationError, errors.Cause(err))
	assert.Nil(t, listItemResult)
}

func TestService_GetItemByIDSuccess(t *testing.T) {
	itemExpected := Item{
		ID:          1,
		Name:        "test",
		Image:       "test",
		Price:       10000,
		Description: "test",
	}
	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemByID", 10, 1).Return(itemExpected, nil)

	// Test
	itemResult, err := mockService.GetItemByID(10, 1)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &itemExpected, itemResult)
	assert.NotNil(t, itemResult)
	assert.NoError(t, err)
}

func TestService_GetItemByIDError(t *testing.T) {
	item := Item{}
	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo, nil)

	mockRepo.On("GetItemByID", 10, 1).Return(item, ErrInternalServerError)

	// Test
	itemResult, err := mockService.GetItemByID(10, 1)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, itemResult)
}

func TestService_DeleteItemAdminByID(t *testing.T) {
	t.Run("success status completed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		itemID := 1

		mockRepo.On("DeleteItemAdminByID", itemID).Return(nil)

		err := service.DeleteItemAdminByID(itemID)
		assert.Nil(t, err)
	})

	t.Run("failed status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo, nil)

		// input
		itemID := 1

		mockRepo.On("DeleteItemAdminByID", itemID).Return(ErrInternalServerError)

		err := service.DeleteItemAdminByID(itemID)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}

func TestService_UpdateItem(t *testing.T) {
	imageString := "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg=="

	tests := map[string]struct {
		expectedItem Item
		wantError    error
		cloudinary   cloudinary.Repo
	}{
		"success": {
			expectedItem: Item{
				Name:        "Nama item",
				Image:       imageString,
				Description: "Deskripsi item",
				Price:       1000,
			},
			wantError: nil,
		},
		"image string is not data uri": {
			expectedItem: Item{
				Image: "blablablas:image/jpeg;base64,==",
			},
			wantError: ErrInputValidationError,
		},
		"image string is not image": {
			expectedItem: Item{
				Image: "data:not-image/jpeg;base64,==",
			},
			wantError: ErrInputValidationError,
		},
		"image string is not base64 encoded": {
			expectedItem: Item{
				Image: "data:image/jpeg;base64,==",
			},
			wantError: ErrInputValidationError,
		},
		"failed to upload image to cloudinary": {
			expectedItem: Item{
				Name:        "Nama item",
				Image:       "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWBAMAAADOL2zRAAAAG1BMVEXMzMyWlpaqqqq3t7fFxcW+vr6xsbGjo6OcnJyLKnDGAAAACXBIWXMAAA7EAAAOxAGVKw4bAAABAElEQVRoge3SMW+DMBiE4YsxJqMJtHOTITPeOsLQnaodGImEUMZEkZhRUqn92f0MaTubtfeMh/QGHANEREREREREREREtIJJ0xbH299kp8l8FaGtLdTQ19HjofxZlJ0m1+eBKZcikd9PWtXC5DoDotRO04B9YOvFIXmXLy2jEbiqE6Df7DTleA5socLqvEFVxtJyrpZFWz/pHM2CVte0lS8g2eDe6prOyqPglhzROL+Xye4tmT4WvRcQ2/m81p+/rdguOi8Hc5L/8Qk4vhZzy08DduGt9eVQyP2qoTM1zi0/uf4hvBWf5c77e69Gf798y08L7j0RERERERERERH9P99ZpSVRivB/rgAAAABJRU5ErkJggg==",
				Description: "Deskripsi item",
				Price:       1000,
			},
			wantError: ErrInternalServerError,
			cloudinary: MockCloudinary{
				isError: true,
			},
		},
		"internal error from repository": {
			expectedItem: Item{
				Name:        "Nama item",
				Image:       imageString,
				Description: "Deskripsi item",
				Price:       1000,
			},
			wantError: ErrInternalServerError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo, MockCloudinary{
				imageString: imageString,
			})
			expectedID := 1

			mockRepo.On("UpdateItem", expectedID, test.expectedItem).Return(test.wantError)

			err := service.UpdateItem(expectedID, test.expectedItem)
			if test.wantError != nil {
				assert.True(t, errors.Is(err, test.wantError))
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
