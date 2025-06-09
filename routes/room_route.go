package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func Room(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	roomController := do.MustInvoke[controller.RoomController](injector)

	routes := route.Group("/api/room")
	{
		// Room
		routes.POST("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), roomController.Create)
		routes.GET("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), roomController.GetAllRoom)
		routes.GET("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), roomController.GetRoomByID)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), roomController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "departemen"), roomController.Delete)
	}
}
