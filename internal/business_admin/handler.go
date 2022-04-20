package businessadmin

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	userID := user.ID

	balanceDetail, err := h.service.GetBalanceDetail(userID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data:    balanceDetail,
	})
}

// GetListTransactionsHistoryWithPagination is a handler for API request to get list of transaction history in business admin
func (h *Handler) GetListTransactionsHistoryWithPagination(c echo.Context) error {
	errorList := []string{}
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := ListTransactionRequest{}
	params.UserID = user.ID
	params.Limit = limit
	params.Page = page
	params.Path = "/api/v1/business-admin/transaction-history"

	listTransaction, pagination, err := h.service.GetListTransactionsHistoryWithPagination(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"transactions": listTransaction.Transactions,
			"pagination":   pagination,
		},
	})
}

func (h *Handler) GetTransactionHistoryDetail(c echo.Context) error {
	errorList := []string{}

	_, _, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	bookinIDString := c.Param("bookingID")
	bookingID, err := strconv.Atoi(bookinIDString)
	if err != nil {
		errorList = append(errorList, "bookingID must be number")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	transactionHistoryDetail, err := h.service.GetTransactionHistoryDetail(bookingID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data:    transactionHistoryDetail,
	})
}
