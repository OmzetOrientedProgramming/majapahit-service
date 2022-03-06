package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
)

type Routes struct {
	Router         *echo.Echo
	checkUPHandler *checkup.Handler
	placeHandler   *place.Handler
}

func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler, placeHandler *place.Handler) *Routes {
	return &Routes{
		Router:         router,
		checkUPHandler: checkUpHandler,
		placeHandler:   placeHandler,
	}
}

func (r *Routes) Init() {
	// Application check up
	r.Router.GET("/", r.checkUPHandler.GetApplicationCheckUp)

	// V1
	v1 := r.Router.Group("/api/v1")

	// Place module
	place := v1.Group("/place")
	place.GET("", r.placeHandler.GetPlacesListWithPagination)

}
