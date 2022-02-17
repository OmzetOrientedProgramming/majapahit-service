package checkup

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetApplicationCheckUp() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestGetApplicationCheckUpSuccess(t *testing.T) {
	mockRepo := new(MockRepository)
	testService := NewService(mockRepo)

	// Expectations
	mockRepo.On("GetApplicationCheckUp").Return(true, nil)

	// Testing
	up, err := testService.GetApplicationCheckUp()
	mockRepo.AssertExpectations(t)
	assert.Equal(t, true, up)
	assert.Equal(t, nil, err)
}

func TestService_GetApplicationCheckUpFailed(t *testing.T) {
	mockRepo := new(MockRepository)
	testService := NewService(mockRepo)

	// Mock Repo Function
	mockRepo.On("GetApplicationCheckUp").Return(false, errors.New("database connection error"))

	up, err := testService.GetApplicationCheckUp()
	mockRepo.AssertExpectations(t)
	assert.Equal(t, false, up)
	assert.Error(t, err, "database connection error")
}
