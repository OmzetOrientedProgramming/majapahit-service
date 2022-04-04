package businessadmin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler for defining handler struct
type Handler struct {
	service Service
}

// NewHandler for initialize handler struct
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetBalanceDetail is the handler for getting current balance and latest disbursement date
func (h *Handler) GetBalanceDetail(c echo.Context) error {
	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
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

	userID := user.ID

	balanceDetail, err := h.service.GetBalanceDetail(userID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}
		logrus.Error("[error while calling booking service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data:    balanceDetail,
	})
}
