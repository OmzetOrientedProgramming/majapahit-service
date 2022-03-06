package item

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

func (h *Handler) GetListItem(c echo.Context) error {
	placeIDString := c.Param("placeID")
	name := c.QueryParam("name")

	placeID, _ := strconv.Atoi(placeIDString)

	listItem, _ := h.service.GetListItem(placeID, name)

	return c.JSON(http.StatusOK, util.APIResponse{
		Status:  200,
		Message: "success",
		Data: map[string]interface{}{
			"places":     listItem.Items,
		},
	})
}