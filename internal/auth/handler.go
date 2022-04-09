package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
)

// Handler for handling auth endpoint
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CheckPhoneNumber for handling CheckPhoneNumber endpoint
func (h *Handler) CheckPhoneNumber(c echo.Context) error {
	session := c.QueryParam("session")
	if !(session == "register" || session == "login") {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidation, "session must be register or login")
	}

	var req CheckPhoneNumberRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
	}

	exist, err := h.service.CheckPhoneNumber(req.PhoneNumber)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	if session == "register" {
		if exist {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidation, "phone number already registered")
		}
	} else if session == "login" {
		if !exist {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidation, "phone number has not been registered")
		}
	}

	sessionInfo, err := h.service.SendOTP(req.PhoneNumber, req.RecaptchaToken)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: CheckPhoneNumberResponse{
			SessionInfo: sessionInfo,
		},
	})
}

// VerifyOTP for handling VerifyOTP endpoint
func (h *Handler) VerifyOTP(c echo.Context) error {
	var req VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServer, err.Error())
	}

	resp, err := h.service.VerifyOTP(req.SessionInfo, req.OTP)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
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

// Register for handling Register endpoint
func (h *Handler) Register(c echo.Context) error {
	userData, _, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, errors.Wrap(ErrInternalServer, err.Error()))
	}

	customer, err := h.service.Register(Customer{
		PhoneNumber: userData.Users[0].PhoneNumber,
		Name:        req.FullName,
		LocalID:     userData.Users[0].LocalID,
	})
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		logrus.Error("[error while calling user service] ", err.Error())
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "created",
		Data:    customer,
	})
}
