package provider

import (
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideBookingRequestDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	bookingRequestRepository := repository.NewBookingRequestRepository(db)
	eventRepository := repository.NewEventRepository(db)
	roomRepository := repository.NewRoomRepository(db)

	// Service
	bookingRequestService := service.NewBookingRequestService(bookingRequestRepository, roomRepository, eventRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.BookingRequestController, error) {
			return controller.NewBookingRequestController(bookingRequestService), nil
		},
	)
}
