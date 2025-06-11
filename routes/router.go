package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/samber/do"
)

func BookingRequest(server *gin.Engine, injector *do.Injector) {
	bookingRequestController := do.MustInvokeNamed[controller.BookingRequestController](injector, constants.BookingRequestController)
	BookingRequestRoutes(server, bookingRequestController)
}

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	User(server, injector)
	Department(server, injector)
	Event(server, injector)
	Room(server, injector)
	Invitation(server, injector)
	BookingRequest(server, injector)
}
