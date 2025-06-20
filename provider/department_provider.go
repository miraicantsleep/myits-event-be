package provider

import (
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideDepartmentDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	departmentRepository := repository.NewDepartmentRepository(db)
	userRepository := repository.NewUserRepository(db)

	// Service
	departmentService := service.NewDepartmentService(departmentRepository, userRepository, jwtService, db)
	userService := service.NewUserService(userRepository, jwtService, db)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.DepartmentController, error) {
			return controller.NewDepartmentController(departmentService, userService), nil
		},
	)
}
