package customer

import (
	"github.com/labstack/echo/v4"
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

func (h *Handler) PutEditCustomer(c echo.Context) error {
	return nil
}
