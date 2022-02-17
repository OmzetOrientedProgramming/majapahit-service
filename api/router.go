package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service.git/internal/checkup"
)

type Routes struct {
	Router         *echo.Echo
	checkUPHandler *checkup.Handler
}

func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler) *Routes {
	return &Routes{
		Router:         router,
		checkUPHandler: checkUpHandler,
	}
}

func (r *Routes) Init() {
	// Application check up
	r.Router.GET("/", r.checkUPHandler.GetApplicationCheckUp)
}
