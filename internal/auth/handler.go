package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"time"
)

type Handler struct {
	service   Service
	jwtSecret string
}

// JWTCustomClaims represent custom claims extending default ones.
type JWTCustomClaims struct {
	UserID int `json:"user_id"`
	Status int `json:"status"`
	jwt.StandardClaims
}

// NewHandler is used to initialize Handler
func NewHandler(service Service, JWTSecret string) *Handler {
	return &Handler{
		service:   service,
		jwtSecret: JWTSecret,
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
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})
	}

	switch session {
	case "register":
		if exist {
			return c.JSON(http.StatusConflict, util.APIResponse{
				Status:  http.StatusConflict,
				Message: "phone number already registered",
			})
		}
		if err := h.service.SendOTP(req.PhoneNumber); err != nil {
			logrus.Error("[error while sending otp]", err.Error())
			return c.JSON(http.StatusInternalServerError, util.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: "cannot send otp to phone number",
			})
		}
		return c.JSON(http.StatusOK, util.APIResponse{
			Status:  http.StatusOK,
			Message: "phone number is available",
		})
	case "login":
		if !exist {
			return c.JSON(http.StatusNotFound, util.APIResponse{
				Status:  http.StatusNotFound,
				Message: "phone number has not been registered",
			})
		}
		if err := h.service.SendOTP(req.PhoneNumber); err != nil {
			logrus.Error("[error while sending otp]", err.Error())
			return c.JSON(http.StatusInternalServerError, util.APIResponse{
				Status:  http.StatusInternalServerError,
				Message: "cannot send otp to phone number",
			})
		}
		return c.JSON(http.StatusOK, util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
		})
	}
	return nil
}

// VerifyOTP for handling VerifyOTP endpoint
func (h *Handler) VerifyOTP(c echo.Context) error {
	session := c.QueryParam("session")
	if !(session == "register" || session == "login") {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  []string{"session must be register or login"},
		})
	}

	var req VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		logrus.Error("[error while binding verify otp request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	ok, err := h.service.VerifyOTP(req.PhoneNumber, req.OTP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request to check phone number",
		})
	}

	switch session {
	case "register":
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, util.APIResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "wrong otp code",
			})
		}

		return c.JSON(http.StatusOK, util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
		})
	case "login":
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, util.APIResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "wrong otp code",
			})
		}

		customer, err := h.service.GetCustomerByPhoneNumber(req.PhoneNumber)

		token, err := createJWTToken(1, 1, h.jwtSecret)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, util.APIResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"customer":      *customer,
				"access_token":  token,
				"refresh_token": "not implemented",
			},
		})
	}
	return nil
}

// Register for handling Register endpoint
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		logrus.Error("[error while binding verify otp request]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	customer, err := h.service.Register(Customer{
		PhoneNumber: req.PhoneNumber,
		Name:        req.FullName,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot process request",
		})
	}

	token, err := createJWTToken(customer.ID, customer.Status, h.jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "cannot create JWT token",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"customer":      customer,
			"access_token":  token,
			"refresh_token": "not implemented",
		},
	})
}

func createJWTToken(userID int, status int, secret string) (string, error) {
	claims := &JWTCustomClaims{
		UserID: userID,
		Status: status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
