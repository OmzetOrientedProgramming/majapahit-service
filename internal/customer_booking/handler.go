package customerbooking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

func (h *Handler) GetListCustomerBookingWithPagination(c echo.Context) error {
	placeIDString := c.Param("placeID")
	stateString := c.QueryParam("state")
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")
	
	placeID, _ := strconv.Atoi(placeIDString)
	state, _ := strconv.Atoi(stateString)
	limit, _ := strconv.Atoi(limitString)
	page, _ := strconv.Atoi(pageString)

	params := ListRequest{}
	params.Path = "/api/v1/business-admin/" + placeIDString + "/booking"
	params.PlaceID = placeID
	params.State = state
	params.Limit = limit
	params.Page = page

	listCustomerBooking, pagination, _ := h.service.GetListCustomerBookingWithPagination(params)

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"bookings":     listCustomerBooking.CustomerBookings,
			"pagination": pagination,
		},
	})
}