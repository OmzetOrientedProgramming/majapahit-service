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

func (m *MockRepository) GetListItem(place_id int, name string) (*ListItem, error) {
	args := m.Called(place_id, name)
	ret := args.Get(0).(ListItem)
	return &ret, args.Error(1)
}

func (m *MockRepository) GetItemById(place_id int, item_id int) (*Item, error) {
	args := m.Called(place_id, item_id)
	ret := args.Get(0).(Item)
	return &ret, args.Error(1)
}

func TestService_GetListItemByIDSuccess(t *testing.T) {
	// Define input and output
	listItemExpected := ListItem{
		Items: []Item{
			{
				ID:         	1,
				Name:       	"test 1",
				Image:			"test 1",
				Description:	"test 1",
				Price:     		10000,
			},
			{
				ID:          	2,
				Name:        	"test 2",
				Image:			"test 2",
				Description: 	"test 2",
				Price:     		20000,
			},
		},
	}


	// Init mock repository and mock service
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	// Expectation
	mockRepo.On("GetListItem", 1, "").Return(listItemExpected, nil)

	// Test
	listItemResult, err := mockService.GetListItem(1, "")
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &listItemExpected, listItemResult)
	assert.NotNil(t, listItemResult)
	assert.NoError(t, err)
}

func TestService_GetListItemByIDError(t *testing.T) {
	listItem := ListItem {}
	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetListItem", 1, "").Return(listItem, ErrInternalServerError)

	// Test
	listItemResult, err := mockService.GetListItem(1, "")
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, listItemResult)
}

func TestService_GetItemByIDSuccess(t *testing.T) {
	itemExpected := Item {
		ID: 1,
		Name: "test",
		Image: "test",
		Price: 10000,
		Description: "test",
	}
	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetItemById", 10, 1).Return(itemExpected, nil)

	// Test
	itemResult, err := mockService.GetItemByID(10, 1)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, &itemExpected, itemResult)
	assert.NotNil(t, itemResult)
	assert.NoError(t, err)
}

func TestService_GetItemByIDError(t *testing.T) {
	item := Item {}
	// Mock DB
	mockRepo := new(MockRepository)
	mockService := NewService(mockRepo)

	mockRepo.On("GetItemById", 10, 1).Return(item, ErrInternalServerError)

	// Test
	itemResult, err := mockService.GetItemByID(10, 1)
	mockRepo.AssertExpectations(t)

	assert.Equal(t, ErrInternalServerError, errors.Cause(err))
	assert.Nil(t, itemResult)
}