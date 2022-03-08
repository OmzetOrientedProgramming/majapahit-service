package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
)

type Routes struct {
	Router         *echo.Echo
	checkUPHandler *checkup.Handler
	catalogHandler *item.Handler
}

func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler) *Routes {
	return &Routes{
		Router:         router,
		checkUPHandler: checkUpHandler,
		catalogHandler: catalogHandler,
	}
}

func (r *Routes) Init() {
	// Application check up
	r.Router.GET("/", r.checkUPHandler.GetApplicationCheckUp)

	// V1
	v1 := r.Router.Group("/api/v1")

	// Place module
	place := v1.Group("/place")

	// Catalog module
	catalog := place.Group("/:placeID/catalog")
	catalog.GET("", r.catalogHandler.GetListItem)
	catalog.GET("/:itemID", r.catalogHandler.GetItemByID)
}
