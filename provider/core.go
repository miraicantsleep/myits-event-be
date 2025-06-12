package provider

import (
	"github.com/miraicantsleep/myits-event-be/config"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
	"gorm.io/gorm"
)

func InitDatabase(injector *do.Injector) {
	do.ProvideNamed(injector, constants.DB, func(i *do.Injector) (*gorm.DB, error) {
		return config.SetUpDatabaseConnection(), nil
	})
}

func RegisterDependencies(injector *do.Injector) {
	InitDatabase(injector)

	do.ProvideNamed(injector, constants.JWTService, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	// Initialize
	db := do.MustInvokeNamed[*gorm.DB](injector, constants.DB)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)

	// Provide Dependencies
	ProvideUserDependencies(injector, db, jwtService)
	ProvideDepartmentDependencies(injector, db, jwtService)
	ProvideEventDependencies(injector, db, jwtService)
	ProvideRoomDependencies(injector, db, jwtService)
	ProvideInvitationDependencies(injector, db, jwtService)
	ProvideBookingRequestDependencies(injector, db, jwtService)
}
