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
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"session must be register or login"},
		})
	}

	var req CheckPhoneNumberRequest
	if err := c.Bind(&req); err != nil {
		logrus.Error("[error while binding phone number request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	exist, err := h.service.CheckPhoneNumber(req.PhoneNumber)
	if err != nil {
		logrus.Error("[error while binding phone number request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})
	}

	if session == "register" {
		if exist {
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: "input validation error",
				Errors:  []string{"phone number already registered"},
			})
		}
		sessionInfo, err := h.service.SendOTP(req.PhoneNumber, req.RecaptchaToken)
		if err != nil {
			if errors.Cause(err) == ErrInputValidation {
				errList, errMessage := util.ErrorUnwrap(err)
				return c.JSON(http.StatusBadRequest, util.APIResponse{
					Status:  http.StatusBadRequest,
					Message: errMessage,
					Errors:  errList,
				})
			}
			logrus.Error("[error while sending otp]", err.Error())
			return c.JSON(http.StatusInternalServerError, util.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: "cannot send otp to phone number",
			})
		}
		return c.JSON(http.StatusOK, util.APIResponse{
			Status:  http.StatusOK,
			Message: "phone number is available",
			Data: CheckPhoneNumberResponse{
				SessionInfo: sessionInfo,
			},
		})
	}

	if !exist {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"phone number has not been registered"},
		})
	}
	sessionInfo, err := h.service.SendOTP(req.PhoneNumber, req.RecaptchaToken)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}
		logrus.Error("[error while sending otp]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot send otp to phone number",
		})
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
		logrus.Error("[error while binding verify otp request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	resp, err := h.service.VerifyOTP(req.SessionInfo, req.OTP)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}
		logrus.Error("[error while binding verify otp request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    resp,
	})
}

// Register for handling Register endpoint
func (h *Handler) Register(c echo.Context) error {
	userData, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			errs, message := util.ErrorUnwrap(err)
			return c.JSON(http.StatusForbidden, util.APIResponse{
				Status:  http.StatusForbidden,
				Message: message,
				Errors:  errs,
			})
		}
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		logrus.Error("[error while binding verify otp request] ", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	customer, err := h.service.Register(Customer{
		PhoneNumber: userData.Users[0].PhoneNumber,
		Name:        req.FullName,
		LocalID:     userData.Users[0].LocalID,
	})
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			errs, message := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: message,
				Errors:  errs,
			})
		}
		logrus.Error("[error while calling user service] ", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "created",
		Data:    customer,
	})
}
