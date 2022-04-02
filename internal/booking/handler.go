package booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler struct for booking package
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

	if len(errorList) != 0 {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	var req UpdateBookingStatusRequest
	if err = c.Bind(&req); err != nil {
		logrus.Error("[error while binding update booking status request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	err = h.service.UpdateBookingStatus(bookingID, req.Status)

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
		Message: "Success update status",
	})
}

// GetMyBookingsOngoing will retrieve information related to a customer booking history
func (h *Handler) GetMyBookingsOngoing(c echo.Context) error {
	errorList := []string{}
	userData, err := middleware.ParseUserData(c, util.StatusCustomer)
	localID := userData.Users[0].LocalID

	if len(errorList) != 0 {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	myBookingsOngoing, err := h.service.GetMyBookingsOngoing(localID)
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
		Data:    myBookingsOngoing,
	})
}

// GetMyBookingsPreviousWithPagination will be used to handling the API request for get previous my bookings
func (h *Handler) GetMyBookingsPreviousWithPagination(c echo.Context) error {
	errorList := []string{}
	userData, err := middleware.ParseUserData(c, util.StatusCustomer)
	localID := userData.Users[0].LocalID
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		if limitString == "" {
			limit = 0
		} else {
			errorList = append(errorList, "limit should be positive integer")
		}
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		if pageString == "" {
			page = 0
		} else {
			errorList = append(errorList, "page should be positive integer")
		}
	}

	if len(errorList) != 0 {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	params := BookingsListRequest{}
	params.Path = "/api/v1/booking/previous"
	params.Limit = limit
	params.Page = page

	myBookingsOngoing, pagination, err := h.service.GetMyBookingsPreviousWithPagination(localID, params)

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
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":   myBookingsOngoing.Bookings,
			"pagination": pagination,
		},
	})
}
