package item

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// DeleteItemAdminByID is a handler for API request for delete item by admin, currently updating is_active column
func (h *Handler) DeleteItemAdminByID(c echo.Context) error {
	errorList := []string{}
	itemIDString := c.Param("itemID")

	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		errorList = append(errorList, "incorrect item id")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	err = h.service.DeleteItemAdminByID(itemID)
	if err != nil {
		return util.ErrorWrapWithContext(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
	})
}

// UpdateItem is a handler for updating item API request by business admin
func (h *Handler) UpdateItem(c echo.Context) error {
	var errorList []string
	itemIDString := c.Param("itemID")
	itemID, err := strconv.Atoi(itemIDString)
	if err != nil {
		errorList = append(errorList, "ID harus berupa angka")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, fmt.Errorf("Request tidak valid"), errorList...)
	}

	var itemRequest UpdateItemRequest
	if err := c.Bind(&itemRequest); err != nil {
		logrus.Errorf("failed to parse request: %v", err)
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, fmt.Errorf("Request tidak valid"))
	}

	if err := h.service.UpdateItem(itemID, Item{
		Name:        itemRequest.Name,
		Image:       itemRequest.Image,
		Description: itemRequest.Description,
		Price:       itemRequest.Price,
	}); err != nil {
		switch {
		case errors.Is(err, ErrInternalServerError):
			return util.ErrorWrapWithContext(c, http.StatusInternalServerError, fmt.Errorf("Terdapat kesalahan pada server"))
		case errors.Is(err, ErrNotFound):
			return util.ErrorWrapWithContext(c, http.StatusNotFound, fmt.Errorf("Item tidak ditemukan"))
		case errors.Is(err, ErrInputValidationError):
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, fmt.Errorf("Request tidak valid"))
		}
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "Berhasil",
	})
}

// CreateItem is a handler for creating item API request by business admin
func (h *Handler) CreateItem(c echo.Context) error {
	_, user, err := middleware.ParseUserData(c, util.StatusBusinessAdmin)
	if err != nil {
		if errors.Cause(err) == middleware.ErrForbidden {
			return util.ErrorWrapWithContext(c, http.StatusForbidden, err)
		}
	}

	var itemRequest Item
	if err := c.Bind(&itemRequest); err != nil {
		logrus.Errorf("failed to parse request: %v", err)
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, fmt.Errorf("Request tidak valid"))
	}

	if err := h.service.CreateItem(user.ID, Item{
		Name:        itemRequest.Name,
		Image:       itemRequest.Image,
		Description: itemRequest.Description,
		Price:       itemRequest.Price,
	}); err != nil {
		switch {
		case errors.Is(err, ErrInputValidationError):
			return util.ErrorWrapWithContext(c, http.StatusBadRequest, fmt.Errorf("Request tidak valid"))
		default:
			return util.ErrorWrapWithContext(c, http.StatusInternalServerError, fmt.Errorf("Terdapat kesalahan pada server"))
		}
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "Berhasil",
	})
}
