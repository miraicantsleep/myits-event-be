package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service" // Not directly needed here, but good for context
	// "gorm.io/gorm" // Not directly needed here
)

func BookingRequestRoutes(router *gin.Engine, bookingRequestController controller.BookingRequestController, jwtService service.JWTService /*db *gorm.DB*/) {
	// authMiddleware := middleware.Authenticate(jwtService, db) // Assuming middleware.Authenticate exists
	bookingRequestRoutes := router.Group("/api/booking-requests")
	{
		// bookingRequestRoutes.Use(authMiddleware) // Apply auth middleware to all booking request routes if needed
		bookingRequestRoutes.POST("/", bookingRequestController.Create)
		bookingRequestRoutes.GET("/", bookingRequestController.GetAll)
		bookingRequestRoutes.GET("/:id", bookingRequestController.GetByID)

		// Routes for approval and rejection - might need specific authorization (e.g., admin only)
		// One way to handle this is to have a separate group with different middleware
		// adminBookingRequestRoutes := router.Group("/api/admin/booking-requests")
		// adminBookingRequestRoutes.Use(authMiddleware, middleware.AuthorizeRole("admin")) // Example for admin-only
		// {
		// 	adminBookingRequestRoutes.PATCH("/:id/approve", bookingRequestController.Approve)
		// 	adminBookingRequestRoutes.PATCH("/:id/reject", bookingRequestController.Reject)
		// }
		// For simplicity now, adding them directly. Consider authorization middleware for these.
		bookingRequestRoutes.PATCH("/:id/approve", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), bookingRequestController.Approve)
		bookingRequestRoutes.PATCH("/:id/reject", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), bookingRequestController.Reject)

	}
}
