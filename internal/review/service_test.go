package review

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) InsertBookingReview(review BookingReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockRepository) RetrievePlaceID(bookingID int) (*int, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockRepository) CheckBookingStatus(bookingID int) (bool, error) {
	args := m.Called(bookingID)
	if args.Get(1) != nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRepository) UpdateBookingStatus(bookingID int) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func TestService_InsertBookingReview(t *testing.T) {
	t.Run("Insert booking review done successfully", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		mockRepo.On("RetrievePlaceID", review.BookingID).Return(placeID, nil)
		mockRepo.On("InsertBookingReview", review).Return(nil)
		mockRepo.On("UpdateBookingStatus", review.BookingID).Return(nil)
		err := mockService.InsertBookingReview(review)

		mockRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("Internal server error from CheckBookingStatus", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 100
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(false, ErrInternalServer)
		err := mockService.InsertBookingReview(review)

		mockRepo.AssertExpectations(t)
		assert.Error(t, ErrInternalServer, err)
	})

	t.Run("Booking is not eligible to be reviewed", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(false, nil)
		err := mockService.InsertBookingReview(review)

		mockRepo.AssertExpectations(t)
		assert.Error(t, ErrInputValidation, err)
	})

	t.Run("Exceeded review content input validation", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := strings.Repeat("Test Review", 100)
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Minimum rating value input validation", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 0

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Maximum rating value input validation", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 6

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, errors.Cause(err), ErrInputValidation)
	})

	t.Run("Internal server error from RetrievePlaceID", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		mockRepo.On("RetrievePlaceID", review.BookingID).Return(nil, ErrInternalServer)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})

	t.Run("Internal server error from InsertBookingReview", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		mockRepo.On("RetrievePlaceID", review.BookingID).Return(placeID, nil)
		mockRepo.On("InsertBookingReview", review).Return(ErrInternalServer)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})

	t.Run("Internal server error from UpdateBookingStatus", func(t *testing.T) {
		var placeID *int = new(int)
		*placeID = 1
		userID := 1
		bookingID := 1
		content := ""
		rating := 5

		review := BookingReview{
			UserID:    userID,
			PlaceID:   *placeID,
			BookingID: bookingID,
			Content:   content,
			Rating:    rating,
		}

		mockRepo := new(MockRepository)
		mockService := NewService(mockRepo)
		mockRepo.On("CheckBookingStatus", review.BookingID).Return(true, nil)
		mockRepo.On("RetrievePlaceID", review.BookingID).Return(placeID, nil)
		mockRepo.On("InsertBookingReview", review).Return(nil)
		mockRepo.On("UpdateBookingStatus", review.BookingID).Return(ErrInternalServer)
		err := mockService.InsertBookingReview(review)

		assert.Equal(t, ErrInternalServer, errors.Cause(err))
	})
}
