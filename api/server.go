package api

import (
	"net/http"
	"os"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/database/postgres"
	businessadminauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
)

// Server struct to for the server dependency
type Server struct {
	Router *echo.Echo
}

// NewServer is used to initialize server
func NewServer(router *echo.Echo) *Server {
	return &Server{
		Router: router,
	}
}

var (
	checkUpRepo    checkup.Repo
	checkUpService checkup.Service
	checkupHandler *checkup.Handler

	catalogRepo    item.Repo
	catalogService item.Service
	catalogHandler *item.Handler

	placeRepo    place.Repo
	placeService place.Service
	placeHandler *place.Handler

	authRepo    auth.Repo
	authService auth.Service
	authHandler *auth.Handler

	businessadminauthRepo    businessadminauth.Repo
	businessadminauthService businessadminauth.Service
	businessadminauthHandler *businessadminauth.Handler
)

// Init all dependency
func (s Server) Init() {
	// Init DB
	db := postgres.Init()

	// Init internal module
	// Check up module
	checkUpRepo = checkup.NewRepo(db)
	checkUpService = checkup.NewService(checkUpRepo)
	checkupHandler = checkup.NewHandler(checkUpService)

	// Catalog Module
	catalogRepo = item.NewRepo(db)
	catalogService = item.NewService(catalogRepo)
	catalogHandler = item.NewHandler(catalogService)

	// Place module
	placeRepo = place.NewRepo(db)
	placeService = place.NewService(placeRepo)
	placeHandler = place.NewHandler(placeService)

	// BusinessAdminAuth module
	businessadminauthRepo = businessadminauth.NewRepo(db)
	businessadminauthService = businessadminauth.NewService(businessadminauthRepo)
	businessadminauthHandler = businessadminauth.NewHandler(businessadminauthService)

	// Check up module
	checkUpRepo = checkup.NewRepo(db)
	checkUpService = checkup.NewService(checkUpRepo)
	checkupHandler = checkup.NewHandler(checkUpService)

	// Auth module
	authRepo = auth.NewRepo(db)
	authService = auth.NewService(authRepo, auth.TwillioCredentials{
		SID:        os.Getenv("TWILIO_SID"),
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	authHandler = auth.NewHandler(authService, os.Getenv("JWT_SECRET"))

	// Start routing
	r := NewRoutes(s.Router, checkupHandler, catalogHandler, placeHandler, authHandler, businessadminauthHandler)
	r.Init()
}

// RunServer to run the server
func (s Server) RunServer(port string) {
	if err := s.Router.Start(":" + port); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}
