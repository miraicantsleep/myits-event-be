package provider

import (
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideEventDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	do.ProvideNamed(injector, constants.EventRepository, func(i *do.Injector) (repository.EventRepository, error) {
		return repository.NewEventRepository(db), nil
	})

	// Service
	do.ProvideNamed(injector, constants.EventService, func(i *do.Injector) (service.EventService, error) {
		eventRepo := do.MustInvokeNamed[repository.EventRepository](i, constants.EventRepository)
		// jwtService is available in the ProvideEventDependencies function's scope
		return service.NewEventService(eventRepo, jwtService, db), nil
	})

	// Controller
	do.Provide(injector, func(i *do.Injector) (controller.EventController, error) {
		eventSvc := do.MustInvokeNamed[service.EventService](i, constants.EventService)
		return controller.NewEventController(eventSvc), nil
	})
}
