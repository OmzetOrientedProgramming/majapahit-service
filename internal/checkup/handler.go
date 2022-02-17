package checkup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetApplicationCheckUp(c echo.Context) error {
	_, err := h.service.GetApplicationCheckUp()
	if err != nil {
		logrus.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, CheckUpAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, CheckUpAPIResponse{
		Status:  http.StatusOK,
		Message: "application up!",
	})
}
