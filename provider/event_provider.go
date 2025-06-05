package provider

import (
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideEventDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	eventRepository := repository.NewEventRepository(db)

	// Service
	eventService := service.NewEventService(eventRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.EventController, error) {
			return controller.NewEventController(eventService), nil
		},
	)
}
