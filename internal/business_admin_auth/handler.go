package businessadminauth

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		logrus.Error("Error while binding register business admin request", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "bad request",
			Errors:  []string{err.Error()},
		})
	}

	result, err := h.service.RegisterBusinessAdmin(request)

	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			logrus.Error("Error while validating request", err.Error())
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: "bad request",
				Errors:  []string{err.Error()},
			})
		}

		// Else
		logrus.Error("Error while using user service", err)
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to register business admin",
		})

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
		logrus.Error("Error while binding login request", err.Error())
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "Terjadi kesalahan dalam memproses permintaan anda",
		})
	}

	accessToken, refreshToken, err := h.service.Login(request.Email, request.Password, request.CaptchaResponse)
	if err != nil {
		if errors.Is(err, ErrInputValidationError) {
			return c.JSON(http.StatusUnprocessableEntity, util.APIResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "Kredensial yang anda berikan tidak valid",
			})
		} else if errors.Is(err, ErrUnauthorized) {
			return c.JSON(http.StatusUnauthorized, util.APIResponse{
				Status:  http.StatusUnauthorized,
				Message: "Kredensial yang anda berikan salah",
			})
		} else if errors.Is(err, ErrNotFound) {
			return c.JSON(http.StatusNotFound, util.APIResponse{
				Status:  http.StatusNotFound,
				Message: "Akun tidak ditemukan",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, util.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: "Terjadi kesalahan dalam memproses permintaan anda",
			})
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
