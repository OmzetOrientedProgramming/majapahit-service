package place

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
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

func (h *Handler) GetPlaceDetail(c echo.Context) error {
	placeIdString := c.Param("placeId")

	placeId, err := strconv.Atoi(placeIdString)
	if err != nil {

	}

	placeDetail, err := h.service.GetPlaceDetail(placeId)

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data:    placeDetail,
	})
}

// GetPlacesListWithPagination will be used to handling the API request for get places
func (h *Handler) GetPlacesListWithPagination(c echo.Context) error {
	errorList := []string{}
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		if limitString == "" {
			limit = 0
		} else {
			errorList = append(errorList, "limit should be positive integer")
		}
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		if pageString == "" {
			page = 0
		} else {
			errorList = append(errorList, "page should be positive integer")
		}
	}

	if len(errorList) != 0 {
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	params := PlacesListRequest{}
	params.Path = "/api/v1/place"
	params.Limit = limit
	params.Page = page

	placesList, pagination, err := h.service.GetPlaceListWithPagination(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}

		logrus.Error("[error while accessing place service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
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
