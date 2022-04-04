package api

import (
	"net/http"
	"os"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/booking"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/client"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/database/postgres"
	businessadmin "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin"
	businessadminauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/user"
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

	firebaseAuthRepo firebaseauth.Repo
	authMiddleware   middleware.AuthMiddleware
	authRepo         auth.Repo
	authService      auth.Service
	authHandler      *auth.Handler

	businessadminauthRepo    businessadminauth.Repo
	businessadminauthService businessadminauth.Service
	businessadminauthHandler *businessadminauth.Handler

	bookingRepo    booking.Repo
	bookingService booking.Service
	bookingHandler *booking.Handler

	businessadminRepo    businessadmin.Repo
	businessadminService businessadmin.Service
	businessadminHandler *businessadmin.Handler

	userRepo user.Repo
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
	userRepo = user.NewRepo(db)
	authRepo = auth.NewRepo(db)
	firebaseAuthRepo = firebaseauth.NewRepo(os.Getenv("IDENTITY_TOOLKIT_URL"), os.Getenv("SECURE_TOKEN_URL"), os.Getenv("FIREBASE_API_KEY"))
	authMiddleware = middleware.NewAuthMiddleware(firebaseAuthRepo, userRepo)
	authService = auth.NewService(authRepo, firebaseAuthRepo)
	authHandler = auth.NewHandler(authService)

	// Booking Module
	bookingRepo = booking.NewRepo(db)
	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	xenditService := xendit.NewXenditClient(xenCli)
	bookingService = booking.NewService(bookingRepo, xenditService)
	bookingHandler = booking.NewHandler(bookingService)

	// BusinessAdmin module
	businessadminRepo = businessadmin.NewRepo(db)
	businessadminService = businessadmin.NewService(businessadminRepo)
	businessadminHandler = businessadmin.NewHandler(businessadminService)

	// Start routing
	r := NewRoutes(s.Router, checkupHandler, catalogHandler, placeHandler, authHandler, businessadminauthHandler, authMiddleware, bookingHandler, businessadminHandler)
	r.Init()
}

// RunServer to run the server
func (s Server) RunServer(port string) {
	if err := s.Router.Start(":" + port); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}
