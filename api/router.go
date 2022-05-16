package api

import (
	"github.com/labstack/echo/v4"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/booking"
	businessadmin "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin"
	businessadminauth "gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/business_admin_auth"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/checkup"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/customer"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/item"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/place"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/internal/upload"
	"gitlab.cs.ui.ac.id/ppl-fasilkom-ui/2022/Kelas-B/OOP/majapahit-service/middleware"
)

// Routes struct for routing endpoint
type Routes struct {
	Router                   *echo.Echo
	checkUPHandler           *checkup.Handler
	itemHandler              *item.Handler
	placeHandler             *place.Handler
	authHandler              *auth.Handler
	businessadminauthHandler *businessadminauth.Handler
	authMiddleware           middleware.AuthMiddleware
	bookingHandler           *booking.Handler
	businessadminHandler     *businessadmin.Handler
	customerHandler          *customer.Handler
	uploadHandler            *upload.Handler
}

// NewRoutes for creating Routes instance
func NewRoutes(router *echo.Echo, checkUpHandler *checkup.Handler, itemHandler *item.Handler, placeHandler *place.Handler, authHandler *auth.Handler, businessadminauthHandler *businessadminauth.Handler, authMiddleware middleware.AuthMiddleware, bookingHandler *booking.Handler, businessadminHandler *businessadmin.Handler, customerHandler *customer.Handler, uploadHandler *upload.Handler) *Routes {
	return &Routes{
		Router:                   router,
		checkUPHandler:           checkUpHandler,
		authHandler:              authHandler,
		itemHandler:              itemHandler,
		placeHandler:             placeHandler,
		businessadminauthHandler: businessadminauthHandler,
		authMiddleware:           authMiddleware,
		bookingHandler:           bookingHandler,
		businessadminHandler:     businessadminHandler,
		customerHandler:          customerHandler,
		uploadHandler:            uploadHandler,
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
			catalogRoutes.GET("", r.itemHandler.GetListItemWithPagination)
			catalogRoutes.GET("/:itemID", r.itemHandler.GetItemByID)

			placeRoutes.GET("/:placeID/time-slot", r.bookingHandler.GetTimeSlots, r.authMiddleware.AuthMiddleware())
			placeRoutes.GET("/:placeID/review", r.placeHandler.GetListReviewAndRatingWithPagination)
		}

		// Business Admin Module
		businessAdminRoutes := v1.Group("/business-admin", r.authMiddleware.AuthMiddleware())
		businessAdminRoutes.GET("/balance", r.businessadminHandler.GetBalanceDetail)
		{
			businessAdminRoutes.POST("/disbursement", r.businessadminHandler.CreateDisbursement)

			// Booking Module
			bookingRoutes := businessAdminRoutes.Group("/booking")
			bookingRoutes.GET("", r.bookingHandler.GetListCustomerBookingWithPagination)
			bookingRoutes.GET("/:bookingID", r.bookingHandler.GetDetail)
			bookingRoutes.PATCH("/:bookingID/confirmation", r.bookingHandler.UpdateBookingStatus)

			// List Items Module
			businessProfileRoutes := businessAdminRoutes.Group("/business-profile")
			listItemsRoutes := businessProfileRoutes.Group("/list-items")
			listItemsRoutes.GET("", r.itemHandler.GetListItemAdminWithPagination)
			listItemsRoutes.DELETE("/:itemID", r.itemHandler.DeleteItemAdminByID)

			transactionHistoryRoutes := businessAdminRoutes.Group("/transaction-history")
			transactionHistoryRoutes.GET("", r.businessadminHandler.GetListTransactionsHistoryWithPagination)
			transactionHistoryRoutes.GET("/:bookingID", r.businessadminHandler.GetTransactionHistoryDetail)
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
			bookingRoutes.GET("/detail/:bookingID", r.bookingHandler.GetDetailBookingSaya)
		}

		// callback
		callbackRoutes := v1.Group("/callback")
		{
			xenditCallbackRoutes := callbackRoutes.Group("/xendit")
			{
				xenditCallbackRoutes.POST("/invoices", r.bookingHandler.XenditInvoicesCallback)
				xenditCallbackRoutes.POST("/disbursement", r.businessadminHandler.XenditDisbursementCallback)
			}
		}

		// Customer module
		customerRoutes := v1.Group("/user", r.authMiddleware.AuthMiddleware())
		{
			customerRoutes.PUT("", r.customerHandler.PutEditCustomer)
			customerRoutes.GET("", r.customerHandler.RetrieveCustomerProfile)
		}

		// Upload module
		uploadRoutes := v1.Group("/upload", r.authMiddleware.AuthMiddleware())
		{
			uploadRoutes.POST("/profile-picture", r.uploadHandler.UploadProfilePicture)
		}
	}
}
