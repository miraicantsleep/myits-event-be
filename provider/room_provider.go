package provider

import (
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/repository"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideRoomDependencies(injector *do.Injector, db *gorm.DB, jwtService service.JWTService) {
	// Repository
	roomRepository := repository.NewRoomRepository(db)
	departmentRepository := repository.NewDepartmentRepository(db)

	// Service
	roomService := service.NewRoomService(roomRepository, jwtService, db, departmentRepository)

	// Controller
	do.Provide(
		injector, func(i *do.Injector) (controller.RoomController, error) {
			return controller.NewRoomController(roomService), nil
		},
	)
	do.ProvideNamed(injector, constants.RoomRepository, func(injector *do.Injector) (repository.RoomRepository, error) {
		// Ensure the DB is available
		db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
		return repository.NewRoomRepository(db), nil
	})

}
