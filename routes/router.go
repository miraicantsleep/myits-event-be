package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func BookingRequest(server *gin.Engine, injector *do.Injector) {
	bookingRequestController := do.MustInvokeNamed[controller.BookingRequestController](injector, constants.BookingRequestController)
	jwtService := do.MustInvoke[service.JWTService](injector)
	BookingRequestRoutes(server, bookingRequestController, jwtService)
}

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	User(server, injector)
	Department(server, injector)
	Event(server, injector)
	Room(server, injector)
	Invitation(server, injector)
	BookingRequest(server, injector)
}
