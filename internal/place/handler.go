package place

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
