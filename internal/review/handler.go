package review

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler struct for customer package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// InsertBookingReview for handling booking review endpoint
func (h *Handler) InsertBookingReview(c echo.Context) error {
	_, userModel, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var review BookingReview
	err = c.Bind(&review)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, errors.Wrap(ErrInternalServer, err.Error()))
	}

	review.UserID = userModel.ID

	err = h.service.InsertBookingReview(review)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "Booking review is successfully recorded.",
	})

}