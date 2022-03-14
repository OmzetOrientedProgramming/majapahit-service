package checkup

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Handler struct for checkup package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetApplicationCheckUp for handling GetApplicationCheckUp endpoint
func (h *Handler) GetApplicationCheckUp(c echo.Context) error {
	_, err := h.service.GetApplicationCheckUp()
	if err != nil {
		logrus.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "application up!",
	})
}
