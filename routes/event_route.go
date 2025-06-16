package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func Event(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	eventController := do.MustInvoke[controller.EventController](injector)

	routes := route.Group("/api/event")
	{
		// Event
		routes.GET("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin", "departemen"), eventController.GetAllEvent)
		routes.GET("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), eventController.GetEventByID)
		routes.GET("/:id/attendees", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), eventController.GetEventAttendees)
		routes.POST("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), eventController.Create)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), eventController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa", "admin"), eventController.Delete)
	}
}
