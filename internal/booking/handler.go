package booking

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
	bookingIDString := c.Param("bookingID")
	bookingID, err := strconv.Atoi(bookingIDString)
	if err != nil {
		panic("error atoi")
	}

	bookingDetail, err := h.service.GetDetail(bookingID)
	if err != nil {
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
