package upload

import (
	"github.com/labstack/echo/v4"
)

// Handler struct for upload
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) UploadProfilePicture(c echo.Context) error {
	return nil
}
