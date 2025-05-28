package provider

import (
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideUserDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	userRepository := repository.NewUserRepository(db)

	// Service
	userService := service.NewUserService(userRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.UserController, error) {
			return controller.NewUserController(userService), nil
		},
	)
}
