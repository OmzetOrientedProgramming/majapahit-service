package booking

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler for defining handler struct
type Handler struct {
	service Service
}

// NewHandler for initialize handler struct
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetListCustomerBookingWithPagination is a handler for API request for get customer bookings
func (h *Handler) GetListCustomerBookingWithPagination(c echo.Context) error {
	errorList := []string{}
	stateString := c.QueryParam("state")
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}
	userID := user.ID

	state, err := strconv.Atoi(stateString)
	if err != nil {
		if stateString == "" {
			state = 0
		} else {
			errorList = append(errorList, "state should be positive integer")
		}
	}

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := ListRequest{}
	params.Path = "/api/v1/business-admin/booking"
	params.UserID = userID
	params.State = state
	params.Limit = limit
	params.Page = page

	listCustomerBooking, pagination, err := h.service.GetListCustomerBookingWithPagination(params)

	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
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

// GetAvailableTime for handling get available time endpoint
func (h Handler) GetAvailableTime(c echo.Context) error {
	_, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	errorList := make([]string, 0)

	countString := c.QueryParam("count")
	count, err := strconv.Atoi(countString)
	if err != nil {
		errorList = append(errorList, "count must be number")
	}

	placeIDString := c.Param("placeID")
	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "placeID must be number")
	}

	dateString := c.QueryParam("date")
	date, err := time.Parse(util.DateLayout, dateString)
	if err != nil {
		errorList = append(errorList, "date must be in YYYY-mm-dd format")
	}

	checkInString := c.QueryParam("check_in")
	checkIn, err := time.Parse(util.TimeLayout, checkInString)
	if err != nil {
		errorList = append(errorList, "check_int must be in HH:mm:ss format")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := GetAvailableTimeParams{
		PlaceID:      placeID,
		SelectedDate: date,
		StartTime:    checkIn,
		BookedSlot:   count,
	}

	resp, err := h.service.GetAvailableTime(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    resp,
	})
}

// GetAvailableDate for handling get available time endpoint
func (h Handler) GetAvailableDate(c echo.Context) error {
	_, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	errorList := make([]string, 0)

	placeIDString := c.Param("placeID")
	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "placeID must be number")
	}

	countString := c.QueryParam("count")
	count, err := strconv.Atoi(countString)
	if err != nil {
		errorList = append(errorList, "count must be number")
	}

	intervalString := c.QueryParam("interval")
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		errorList = append(errorList, "interval must be number")
	}

	dateString := c.QueryParam("date")
	date, err := time.Parse(util.DateLayout, dateString)
	if err != nil {
		errorList = append(errorList, "interval must be in YYYY-mm-dd format")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := GetAvailableDateParams{
		PlaceID:    placeID,
		StartDate:  date,
		Interval:   interval,
		BookedSlot: count,
	}

	resp, err := h.service.GetAvailableDate(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    resp,
	})
}

// CreateBooking for handling create booking endpoint
func (h Handler) CreateBooking(c echo.Context) error {
	_, userFromDatabase, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var errorList []string
	var req CreateBookingRequestBody
	if err = c.Bind(&req); err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServerError, err.Error())
	}

	date, err := time.Parse(util.DateLayout, req.Date)
	if err != nil {
		errorList = append(errorList, "date must be in YYYY-mm-dd format")
	}

	startTime, err := time.Parse(util.TimeLayout, req.StartTime)
	if err != nil {
		errorList = append(errorList, "start_time must be in HH:mm:ss format")
	}

	endTime, _ := time.Parse(util.TimeLayout, req.EndTime)
	if err != nil {
		errorList = append(errorList, "end_time must be in HH:mm:ss format")
	}

	placeIDString := c.Param("placeID")
	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "place id must be a number")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	serviceRequest := CreateBookingServiceRequest{
		Items:               req.Items,
		Date:                date,
		StartTime:           startTime,
		EndTime:             endTime,
		Count:               req.Count,
		PlaceID:             placeID,
		UserID:              userFromDatabase.ID,
		CustomerName:        userFromDatabase.Name,
		CustomerPhoneNumber: userFromDatabase.PhoneNumber,
	}

	resp, err := h.service.CreateBooking(serviceRequest)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    resp,
	})
}

// GetTimeSlots for handling time slots endpoint
func (h Handler) GetTimeSlots(c echo.Context) error {
	_, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var errorList []string
	date, err := time.Parse(util.DateLayout, c.QueryParam("date"))
	if err != nil {
		errorList = append(errorList, "date must be in YYYY-mm-dd format")
	}

	placeIDString := c.Param("placeID")
	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "place id must be a number")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	resp, err := h.service.GetTimeSlots(placeID, date)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	var formattedResp []TimeSlotAPIResponse
	for _, i := range *resp {
		formattedResp = append(formattedResp, TimeSlotAPIResponse{
			ID:        i.ID,
			StartTime: i.StartTime.Format(util.TimeLayout),
			EndTime:   i.EndTime.Format(util.TimeLayout),
			Day:       i.Day,
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    formattedResp,
	})
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
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	bookingDetail, err := h.service.GetDetail(bookingID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
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
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	var req UpdateBookingStatusRequest
	if err = c.Bind(&req); err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServerError, err.Error())
	}

	err = h.service.UpdateBookingStatus(bookingID, req.Status)

	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "Success update status",
	})
}

// GetMyBookingsOngoing will retrieve information related to a customer booking history
func (h *Handler) GetMyBookingsOngoing(c echo.Context) error {
	errorList := []string{}
	userData, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	localID := userData.Users[0].LocalID

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	myBookingsOngoing, err := h.service.GetMyBookingsOngoing(localID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)

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
	userData, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	localID := userData.Users[0].LocalID
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := BookingsListRequest{}
	params.Path = "/api/v1/booking/previous"
	params.Limit = limit
	params.Page = page

	myBookingsOngoing, pagination, err := h.service.GetMyBookingsPreviousWithPagination(localID, params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
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

// XenditInvoicesCallback for handling xendit invoices callback
func (h Handler) XenditInvoicesCallback(c echo.Context) error {
	var params XenditInvoicesCallback

	err := c.Bind(&params)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServerError, err.Error())
	}

	err = h.service.XenditInvoicesCallback(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "success",
	})
}

// GetDetailBookingSaya used for handling request to get detail booking saya
func (h *Handler) GetDetailBookingSaya(c echo.Context) error {
	bookingIDString := c.Param("bookingID")
	fmt.Println("bookingID:", bookingIDString)
	bookingID, err := strconv.Atoi(bookingIDString)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, "bookingID must be number")
	}

	detailBookingSaya, err := h.service.GetDetailBookingSaya(bookingID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data:    detailBookingSaya,
	})
}
