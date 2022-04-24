package item

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
	mockService := NewService(mockRepo)

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
		service := NewService(mockRepo)

		// input
		itemID := 1

		mockRepo.On("DeleteItemAdminByID", itemID).Return(nil)

		err := service.DeleteItemAdminByID(itemID)
		assert.Nil(t, err)
	})

	t.Run("failed status", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		// input
		itemID := 1

		mockRepo.On("DeleteItemAdminByID", itemID).Return(ErrInternalServerError)

		err := service.DeleteItemAdminByID(itemID)
		assert.NotNil(t, err)
		assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	})
}