package customerbooking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler struct for item package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetListCustomerBookingWithPagination is a handler for API request for get customer bookings
func (h *Handler) GetListCustomerBookingWithPagination(c echo.Context) error {
	errorList := []string{}
	placeIDString := c.Param("placeID")
	stateString := c.QueryParam("state")
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "incorrect place id")
	}

	state, err := strconv.Atoi(stateString)
	if err != nil {
		if stateString == "" {
			state = 0
		} else {
			errorList = append(errorList, "state should be positive integer")
		}
	}

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

	params := ListRequest{}
	params.Path = "/api/v1/business-admin/" + placeIDString + "/booking"
	params.PlaceID = placeID
	params.State = state
	params.Limit = limit
	params.Page = page

	listCustomerBooking, pagination, err := h.service.GetListCustomerBookingWithPagination(params)

	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}

		logrus.Error("[error while accessing customer booking service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":   listCustomerBooking.CustomerBookings,
			"pagination": pagination,
		},
	})
}
