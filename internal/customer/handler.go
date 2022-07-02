package customer

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
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

// PutEditCustomer for handling edit customer endpoint
func (h *Handler) PutEditCustomer(c echo.Context) error {
	_, userModel, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var req EditCustomerRequest
	err = c.Bind(&req)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, errors.Wrap(ErrInternalServer, err.Error()))
	}

	dateOfBirth, err := time.Parse(util.DateLayout, req.DateOfBirthString)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, errors.Wrap(ErrInputValidation, "Format date of birth tidak sesuai (YYYY-MM-DD)"))
	}

	req.ID = userModel.ID
	req.DateOfBirth = dateOfBirth

	err = h.service.PutEditCustomer(req)
	if err != nil {
		if errors.Cause(err) == ErrInputValidation {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "Successfully Edited Profile!",
	})

}

// RetrieveCustomerProfile returns GET response for customer profile
func (h *Handler) RetrieveCustomerProfile(c echo.Context) error {
	_, user, err := middleware.ParseUserData(c, util.StatusCustomer)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}
	userID := user.ID

	customerProfile, err := h.service.RetrieveCustomerProfile(userID)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    customerProfile,
	})
}
