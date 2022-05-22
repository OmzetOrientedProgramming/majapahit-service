package review

import (
	"strings"

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
	InsertBookingReview(review BookingReview) error
}

type service struct {
	repo Repo
}

func (s service) InsertBookingReview(review BookingReview) error {
	isEligible, err := s.repo.CheckBookingStatus(review.BookingID)
	if err != nil {
		return err
	}
	if !isEligible {
		return errors.Wrap(ErrInputValidation, "Forbidden entry: Wrong booking status")
	}

	var errorList []string

	if len([]rune(review.Content)) > 500 {
		errorList = append(errorList, "Review melebihi 500 karakter.")
	}

	if review.Rating < util.MinimumRatingValue {
		errorList = append(errorList, "Rating invalid. Minimum rating adalah 1.")
	}

	if review.Rating > util.MaximumRatingValue {
		errorList = append(errorList, "Rating invalid. Maksimum rating adalah 5.")
	}

	if len(errorList) > 0 {
		return errors.Wrap(ErrInputValidation, strings.Join(errorList, ";"))
	}

	placeID, err := s.repo.RetrievePlaceID(review.BookingID)
	if err != nil {
		return err
	}
	review.PlaceID = *placeID	

	err = s.repo.InsertBookingReview(review)
	if err != nil {
		return err
	}

	err = s.repo.UpdateBookingStatus(review.BookingID)
	if err != nil {
		return err
	}

	return nil

}