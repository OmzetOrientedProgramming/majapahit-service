package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
)

// Routes struct for routing endpoint
type Routes struct {
	Router         *echo.Echo
	checkUPHandler *checkup.Handler
	catalogHandler *item.Handler
	placeHandler   *place.Handler
}

// NewRoutes for creating Routes instance
func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler, placeHandler *place.Handler) *Routes {
	return &Routes{
		Router:         router,
		checkUPHandler: checkUpHandler,
		catalogHandler: catalogHandler,
		placeHandler:   placeHandler,
	}
}

// Init to init list of endpoint URL
func (r *Routes) Init() {
	// Application check up
	r.Router.GET("/", r.checkUPHandler.GetApplicationCheckUp)

	// V1
	v1 := r.Router.Group("/api/v1")

	// Place module
	place := v1.Group("/place")
	place.GET("", r.placeHandler.GetPlacesListWithPagination)

	// Catalog module
	catalog := place.Group("/:placeID/catalog")
	catalog.GET("", r.catalogHandler.GetListItemWithPagination)
	catalog.GET("/:itemID", r.catalogHandler.GetItemByID)
}
