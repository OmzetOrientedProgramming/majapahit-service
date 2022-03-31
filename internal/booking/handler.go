package booking

import "github.com/labstack/echo/v4"

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
	panic("Not yet implemented!")
}
