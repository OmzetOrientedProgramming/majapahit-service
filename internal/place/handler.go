package place

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
	"net/http"
	"strconv"
)

// Handler struct for place package
type Handler struct {
	service Service
}

// NewHandler is used to initialize Handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetDetail will retrieve information related to a place
func (h *Handler) GetDetail(c echo.Context) error {
	errorList := []string{}
	placeIDString := c.Param("placeID")

	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "placeID must be number")
	}

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	placeDetail, err := h.service.GetDetail(placeID)
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

// GetPlacesListWithPagination will be used to handling the API request for get places
func (h *Handler) GetPlacesListWithPagination(c echo.Context) error {
	errorList := []string{}
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := PlacesListRequest{}
	params.Path = "/api/v1/place"
	params.Limit = limit
	params.Page = page

	placesList, pagination, err := h.service.GetPlaceListWithPagination(params)
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
			"places":     placesList.Places,
			"pagination": pagination,
		},
	})
}
