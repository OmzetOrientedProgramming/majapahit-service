package api

import (
	"net/http"
	"os"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/cloudinary"
	firebaseauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/firebase_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/pkg/xendit"

	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/booking"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/customer"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/upload"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/review"

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

	itemRepo    item.Repo
	itemService item.Service
	itemHandler *item.Handler

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

	customerRepo    customer.Repo
	customerService customer.Service
	customerHandler *customer.Handler

	cloudinaryRepo cloudinary.Repo
	uploadService  upload.Service
	uploadHandler  *upload.Handler

	reviewRepo 		review.Repo
	reviewService  	review.Service
	reviewHandler  	*review.Handler
)

// Init all dependency
func (s Server) Init() {
	// Init DB
	db := postgres.Init()

	cloudinaryRepo = cloudinary.NewRepo(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))

	// Init internal module
	// Check up module
	checkUpRepo = checkup.NewRepo(db)
	checkUpService = checkup.NewService(checkUpRepo)
	checkupHandler = checkup.NewHandler(checkUpService)

	// Catalog Module
	itemRepo = item.NewRepo(db)
	itemService = item.NewService(itemRepo, cloudinaryRepo)
	itemHandler = item.NewHandler(itemService)

	// Place module
	placeRepo = place.NewRepo(db)
	placeService = place.NewService(placeRepo)
	placeHandler = place.NewHandler(placeService)

	// BusinessAdminAuth module
	businessadminauthRepo = businessadminauth.NewRepo(db)
	businessadminauthService = businessadminauth.NewService(businessadminauthRepo, os.Getenv("FIREBASE_API_KEY"), os.Getenv("IDENTITY_TOOLKIT_URL"))
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

	// Xendit service
	xenCli := client.New(os.Getenv("XENDIT_TOKEN"))
	xenditService := xendit.NewXenditClient(xenCli)

	// Booking Module
	bookingRepo = booking.NewRepo(db)
	bookingService = booking.NewService(bookingRepo, xenditService)
	bookingHandler = booking.NewHandler(bookingService)

	// BusinessAdmin module
	businessadminRepo = businessadmin.NewRepo(db)
	businessadminService = businessadmin.NewService(businessadminRepo, xenditService, placeService)
	businessadminHandler = businessadmin.NewHandler(businessadminService)

	// Customer Module
	customerRepo = customer.NewRepo(db)
	customerService = customer.NewService(customerRepo)
	customerHandler = customer.NewHandler(customerService)

	// Upload module
	uploadService = upload.NewService(cloudinaryRepo)
	uploadHandler = upload.NewHandler(uploadService)

	// Review module
	reviewRepo = review.NewRepo(db)
	reviewService = review.NewService(reviewRepo)
	reviewHandler = review.NewHandler(reviewService)

	// Start routing
	r := NewRoutes(s.Router, checkupHandler, itemHandler, placeHandler, authHandler, businessadminauthHandler, authMiddleware, bookingHandler, businessadminHandler, customerHandler, uploadHandler, reviewHandler)
	r.Init()
}

// RunServer to run the server
func (s Server) RunServer(port string) {
	if err := s.Router.Start(":" + port); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}
