package place

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

// GetListReviewAndRatingWithPagination will be used to handling the API request for get review and rating of a place
func (h *Handler) GetListReviewAndRatingWithPagination(c echo.Context) error {
	errorList := []string{}
	placeIDString := c.Param("placeID")
	limitString := c.QueryParam("limit")
	pageString := c.QueryParam("page")
	latestString := c.QueryParam("latest")
	ratingString := c.QueryParam("rating")

	placeID, err := strconv.Atoi(placeIDString)
	if err != nil {
		errorList = append(errorList, "incorrect place id")
	}

	latest, err := strconv.ParseBool(latestString)
	if err != nil {
		if latestString == "" {
			latest = false
		} else {
			errorList = append(errorList, "latest parameter should be boolean type")
		}
	}

	rating, err := strconv.ParseBool(ratingString)
	if err != nil {
		if ratingString == "" {
			rating = false
		} else {
			errorList = append(errorList, "rating parameter should be boolean type")
		}
	}

	page, limit, errorsFromValidator := util.ValidateParams(pageString, limitString)
	errorList = append(errorList, errorsFromValidator...)

	if len(errorList) != 0 {
		return util.ErrorWrapWithContext(c, http.StatusBadRequest, ErrInputValidationError, errorList...)
	}

	params := ListReviewRequest{}
	params.Path = "/api/v1/place/" + placeIDString + "/review"
	params.Limit = limit
	params.Page = page
	params.Latest = latest
	params.Rating = rating
	params.PlaceID = placeID

	listReview, pagination, err := h.service.GetListReviewAndRatingWithPagination(params)
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
			"reviews":      listReview.Reviews,
			"pagination":   pagination,
			"total_review": listReview.TotalCount,
		},
	})
}
