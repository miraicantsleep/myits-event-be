package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func BookingRequest(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	bookingRequestController := do.MustInvoke[controller.BookingRequestController](injector)

	routes := route.Group("/api/booking-request")
	{
		routes.POST("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa"), bookingRequestController.Create)
		routes.GET("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin", "departemen"), bookingRequestController.GetByID)
		routes.GET("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin", "departemen"), bookingRequestController.GetAll)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), bookingRequestController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), bookingRequestController.Delete)
		routes.PATCH("/:id/approve", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), bookingRequestController.Approve)
		routes.PATCH("/:id/reject", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), bookingRequestController.Reject)
		routes.GET("/with-capacity", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), bookingRequestController.GetAllWithCapacity)
	}
}
