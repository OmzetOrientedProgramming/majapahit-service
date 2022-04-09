package item

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

// Handler struct for item package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetListItemWithPagination is a handler for API request for get list item (catalog)
func (h *Handler) GetListItemWithPagination(c echo.Context) error {
	errorList := []string{}
	placeIDString := c.Param("placeID")
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")
	name := c.QueryParam("name")

	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "incorrect place id")
	}

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := ListItemRequest{}
	params.PlaceID = placeID
	params.Name = name
	params.Limit = limit
	params.Page = page
	params.Path = "/api/v1/place/" + placeIDString + "/catalog"

	listItem, pagination, err := h.service.GetListItemWithPagination(params)
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
			"items":      listItem.Items,
			"info":       listItem.PlaceInfo,
			"pagination": pagination,
		},
	})
}

// GetItemByID is A handler for API request for get detail item
func (h *Handler) GetItemByID(c echo.Context) error {
	errorList := []string{}
	placeIDString := c.Param("placeID")
	itemIDString := c.Param("itemID")

	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "incorrect place id")
	}

	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		errorList = append(errorList, "incorrect item id")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	item, err := h.service.GetItemByID(placeID, itemID)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"item": item,
		},
	})
}

// GetListItemAdminWithPagination is a handler for API request to get list items in business admin
func (h *Handler) GetListItemAdminWithPagination(c echo.Context) error {
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

	params := ListItemRequest{}
	params.UserID = user.ID
	params.PlaceID = 0
	params.Limit = limit
	params.Page = page
	params.Path = "/api/v1/business-admin/business-profile/list-items"

	listItem, pagination, err := h.service.GetListItemWithPagination(params)
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
			"items":      listItem.Items,
			"pagination": pagination,
		},
	})
}
