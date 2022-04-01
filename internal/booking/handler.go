package booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler struct for place package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetDetail will retrieve information related to a booking
func (h *Handler) GetDetail(c echo.Context) error {
	errorList := []string{}

	bookingIDString := c.Param("bookingID")
	bookingID, err := strconv.Atoi(bookingIDString)
	if err != nil {
		errorList = append(errorList, "bookingID must be number")
	}

	if len(errorList) != 0 {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	bookingDetail, err := h.service.GetDetail(bookingID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}

		logrus.Error("[error while accessing booking service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    bookingDetail,
	})
}

// UpdateBookingStatus will update booking status
func (h *Handler) UpdateBookingStatus(c echo.Context) error {
	errorList := []string{}

	bookingIDString := c.Param("bookingID")
	bookingID, err := strconv.Atoi(bookingIDString)
	if err != nil {
		errorList = append(errorList, "bookingID must be number")
	}

	var req UpdateBookingStatusRequest
	if err = c.Bind(&req); err != nil {
		panic(err.Error())
	}

	err = h.service.UpdateBookingStatus(bookingID, req.Status)
	if err != nil {
		panic(err.Error())
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "Success update status",
	})
}
