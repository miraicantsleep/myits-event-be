package provider

import (
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideBookingRequestDependencies(injector *do.Injector, db *gorm.DB) {
	// Provide BookingRequestRepository
	do.ProvideNamed(injector, constants.BookingRequestRepository, func(i *do.Injector) (repository.BookingRequestRepository, error) {
		return repository.NewBookingRequestRepository(db), nil
	})

	// Provide BookingRequestService
	do.ProvideNamed(injector, constants.BookingRequestService, func(i *do.Injector) (service.BookingRequestService, error) {
		bookingRequestRepo := do.MustInvokeNamed[repository.BookingRequestRepository](i, constants.BookingRequestRepository)
		roomRepo := do.MustInvokeNamed[repository.RoomRepository](i, constants.RoomRepository)
		eventRepo := do.MustInvokeNamed[repository.EventRepository](i, constants.EventRepository) // Assumes EventRepository is available
		return service.NewBookingRequestService(bookingRequestRepo, roomRepo, eventRepo, db), nil
	})

	// Provide BookingRequestController
	do.ProvideNamed(injector, constants.BookingRequestController, func(i *do.Injector) (controller.BookingRequestController, error) {
		bookingRequestService := do.MustInvokeNamed[service.BookingRequestService](i, constants.BookingRequestService)
		// jwtService := do.MustInvokeNamed[service.JWTService](i, constants.JWTService) // If auth is needed for controller
		return controller.NewBookingRequestController(bookingRequestService /*, jwtService */), nil
	})
}
