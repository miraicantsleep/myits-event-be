package provider

import (
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideInvitationDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	invitationRepository := repository.NewInvitationRepository(db)

	// Service
	invitationService := service.NewInvitationService(invitationRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.InvitationController, error) {
			return controller.NewInvitationController(invitationService, jwtService), nil
		},
	)
}
