package businessadminauth

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
)

// Handler is a struct to define Handler
type Handler struct {
	service Service
}

// NewHandler is a constructor to get Handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterBusinessAdmin is the main media to bind requests, process it and return its response
func (h *Handler) RegisterBusinessAdmin(c echo.Context) error {
	var request RegisterBusinessAdminRequest
	err := c.Bind(&request)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, err.Error())
	}

	result, err := h.service.RegisterBusinessAdmin(request)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    result,
	})
}

// Login is the main media to bind requests, process it and return its response
func (h *Handler) Login(c echo.Context) error {
	var request LoginRequest
	err := c.Bind(&request)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, err.Error())

	}

	accessToken, refreshToken, err := h.service.Login(request.Email, request.Password, request.CaptchaResponse)
	if err != nil {
		if errors.Is(err, ErrInputValidationError) {
			return util.ErrorWrapWithContext(c, http.StatusUnprocessableEntity, err)
		} else if errors.Is(err, ErrUnauthorized) {
			return util.ErrorWrapWithContext(c, http.StatusUnauthorized, err)
		} else if errors.Is(err, ErrNotFound) {
			return util.ErrorWrapWithContext(c, http.StatusNotFound, err)
		} else {
			return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "Login berhasil",
		Data: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
