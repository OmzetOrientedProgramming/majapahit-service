package item

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/util"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

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

	params := ListItemRequest{}
	params.PlaceID = placeID
	params.Name = name
	params.Limit = limit
	params.Page = page
	params.Path = "/api/v1/place/"+placeIDString+"/catalog"

	listItem, pagination, err := h.service.GetListItemWithPagination(params)
	if err != nil {
		if errors.Cause(err) == ErrInputValidationError {
			errList, errMessage := util.ErrorUnwrap(err)
			return c.JSON(http.StatusBadRequest, util.APIResponse{
				Status:  http.StatusBadRequest,
				Message: errMessage,
				Errors:  errList,
			})
		}

		logrus.Error("[error while accessing catalog service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"items":     listItem.Items,
			"pagination": pagination,
		},
	})
}

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
		return c.JSON(http.StatusBadRequest, util.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "input validation error",
			Errors:  errorList,
		})
	}

	item, err := h.service.GetItemByID(placeID, itemID)
	if err != nil {
		logrus.Error("[error while accessing catalog service]", err.Error())
		return c.JSON(http.StatusInternalServerError, util.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"item":     item,
		},
	})
}