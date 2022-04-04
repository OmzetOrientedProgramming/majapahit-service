package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/booking"
	businessadmin "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin"
	businessadminauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
)

// Routes struct for routing endpoint
type Routes struct {
	Router                   *echo.Echo
	checkUPHandler           *checkup.Handler
	catalogHandler           *item.Handler
	placeHandler             *place.Handler
	authHandler              *auth.Handler
	businessadminauthHandler *businessadminauth.Handler
	authMiddleware           middleware.AuthMiddleware
	bookingHandler           *booking.Handler
	businessadminHandler     *businessadmin.Handler
}

// NewRoutes for creating Routes instance
func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler, catalogHandler *item.Handler, placeHandler *place.Handler, authHandler *auth.Handler, businessadminauthHandler *businessadminauth.Handler, authMiddleware middleware.AuthMiddleware, bookingHandler *booking.Handler, businessadminHandler *businessadmin.Handler) *Routes {
	return &Routes{
		Router:                   router,
		checkUPHandler:           checkUpHandler,
		authHandler:              authHandler,
		catalogHandler:           catalogHandler,
		placeHandler:             placeHandler,
		businessadminauthHandler: businessadminauthHandler,
		authMiddleware:           authMiddleware,
		bookingHandler:           bookingHandler,
		businessadminHandler:     businessadminHandler,
	}
}

// Init to init list of endpoint URL
func (r *Routes) Init() {
	// Application check up
	r.Router.GET("/", r.checkUPHandler.GetApplicationCheckUp)

	// V1
	v1 := r.Router.Group("/api/v1")
	{
		// Place module
		placeRoutes := v1.Group("/place")
		placeRoutes.GET("", r.placeHandler.GetPlacesListWithPagination)
		placeRoutes.GET("/:placeID", r.placeHandler.GetDetail)
		{
			// Catalog module
			catalogRoutes := placeRoutes.Group("/:placeID/catalog")
			catalogRoutes.GET("", r.catalogHandler.GetListItemWithPagination)
			catalogRoutes.GET("/:itemID", r.catalogHandler.GetItemByID)

			placeRoutes.GET("/:placeID/time-slot", r.bookingHandler.GetTimeSlots, r.authMiddleware.AuthMiddleware())
		}

		// Business Admin Module
		businessAdminRoutes := v1.Group("/business-admin", r.authMiddleware.AuthMiddleware())
		businessAdminRoutes.GET("/balance", r.businessadminHandler.GetBalanceDetail)
		{
			// Booking
			bookingRoutes := businessAdminRoutes.Group("/booking")
			bookingRoutes.GET("", r.bookingHandler.GetListCustomerBookingWithPagination)
			bookingRoutes.GET("/:bookingID", r.bookingHandler.GetDetail)
			bookingRoutes.PATCH("/:bookingID/confirmation", r.bookingHandler.UpdateBookingStatus)
		}

		// Auth module
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/check-phone-number", r.authHandler.CheckPhoneNumber)
			authRoutes.POST("/verify-otp", r.authHandler.VerifyOTP)
			authRoutes.POST("/register", r.authHandler.Register, r.authMiddleware.AuthMiddleware())

			authRoutes.POST("/business-admin/register", r.businessadminauthHandler.RegisterBusinessAdmin)
			authRoutes.POST("/business-admin/login", r.businessadminauthHandler.Login)
		}

		// Booking module
		bookingRoutes := v1.Group("/booking", r.authMiddleware.AuthMiddleware())
		{
			bookingRoutes.POST("/:placeID", r.bookingHandler.CreateBooking)
			bookingRoutes.GET("/time/:placeID", r.bookingHandler.GetAvailableTime)
			bookingRoutes.GET("/date/:placeID", r.bookingHandler.GetAvailableDate)
			bookingRoutes.GET("/ongoing", r.bookingHandler.GetMyBookingsOngoing)
			bookingRoutes.GET("/previous", r.bookingHandler.GetMyBookingsPreviousWithPagination)
		}

		// callback
		callbackRoutes := v1.Group("/callback")
		{
			xenditCallbackRoutes := callbackRoutes.Group("/xendit")
			{
				xenditCallbackRoutes.POST("/invoices", r.bookingHandler.XenditInvoicesCallback)
			}
		}
	}
}
