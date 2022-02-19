package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service.git/database/postgres"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service.git/internal/checkup"
)

type Server struct {
	Router *echo.Echo
}

func NewServer(router *echo.Echo) *Server {
	return &Server{
		Router: router,
	}
}

var (
	checkUpRepo    checkup.Repo
	checkUpService checkup.Service
	checkupHandler *checkup.Handler
)

func (s Server) Init() {
	// Init DB
	db := postgres.Init()

	// Init internal module
	// Check up module
	checkUpRepo = checkup.NewRepo(db)
	checkUpService = checkup.NewService(checkUpRepo)
	checkupHandler = checkup.NewHandler(checkUpService)

	// Start routing
	r := NewRoutes(s.Router, checkupHandler)
	r.Init()
}

func (s Server) RunServer(port string) {
	if err := s.Router.Start(":" + port); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}