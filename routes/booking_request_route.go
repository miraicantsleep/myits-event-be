package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
)

func BookingRequestRoutes(router *gin.Engine, bookingRequestController controller.BookingRequestController, jwtService service.JWTService) {
	bookingRequestRoutes := router.Group("/api/booking-requests")
	{
		bookingRequestRoutes.POST("/", bookingRequestController.Create)
		bookingRequestRoutes.GET("/", bookingRequestController.GetAll)
		bookingRequestRoutes.GET("/:id", bookingRequestController.GetByID)

		// Subgroup for routes requiring authentication and specific roles
		protected := bookingRequestRoutes.Group("/")
		protected.Use(middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"))
		{
			protected.PATCH("/:id/approve", bookingRequestController.Approve)
			protected.PATCH("/:id/reject", bookingRequestController.Reject)
		}
	}
}
