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

// CreateDisbursement for handling create disbursement endpoint
func (h *Handler) CreateDisbursement(c echo.Context) error {
	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	userID := user.ID

	var (
		amount      float64
		amountParam map[string]interface{}
		ok          bool
	)

	err = c.Bind(&amountParam)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, errors.Wrap(ErrInputValidationError, "invalid body"), err.Error())
	}

	if amount, ok = amountParam["amount"].(float64); !ok {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, errors.Wrap(ErrInputValidationError, "invalid amount"))
	}

	resp, err := h.service.CreateDisbursement(userID, amount)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{Status: http.StatusOK, Message: "success", Data: resp})
}

// XenditDisbursementCallback for handling xendit disbursement callback
func (h *Handler) XenditDisbursementCallback(c echo.Context) error {
	var params DisbursementCallback

	err := c.Bind(&params)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, ErrInternalServerError, err.Error())
	}

	err = h.service.DisbursementCallbackFromXendit(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, util.APIResponse{
		Status:  http.StatusCreated,
		Message: "success",
	})
}

// GetTransactionHistoryDetail is a handler for API request to get detail transaction history
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

// PutEditProfile is a handler for API request to Update Business Profile
func (h *Handler) PutEditProfile(ctx echo.Context) error {
	_, userModel, err := middleware.ParseUserData(ctx, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(ctx, http.StatusForbidden, err)
		}
	}

	userID := userModel.ID

	var req EditProfileRequest
	err = ctx.Bind(&req)
	if err != nil {
		return util.ErrorWrapWithContext(ctx, http.StatusInternalServerError, errors.Wrap(ErrInternalServer, err.Error()))
	}
	req.UserID = userID

	err = h.service.PutEditProfile(req)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(ctx, http.StatusBadRequest, err)
		}
		return util.ErrorWrapWithContext(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "Successfully Edited Profile!",
	})
}

// GetPlaceDetail will retrieve information related to a place
func (h *Handler) GetPlaceDetail(c echo.Context) error {
	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	userID := user.ID

	placeDetail, err := h.service.GetPlaceDetail(userID)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    placeDetail,
	})
}

// GetListReviewAndRatingWithPagination will retrieve information related to a place
func (h *Handler) GetListReviewAndRatingWithPagination(c echo.Context) error {
	errorList := []string{}
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	userID := user.ID

	params := ListReviewRequest{}
	params.Path = "/api/v1/business-admin/business-profile/review"
	params.Limit = limit
	params.Page = page

	listReview, pagination, err := h.service.GetListReviewAndRatingWithPagination(userID, params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, err)
		}

		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"reviews":      listReview.Reviews,
			"pagination":   pagination,
			"total_review": listReview.TotalCount,
		},
	})
}
