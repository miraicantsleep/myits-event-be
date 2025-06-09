package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func Invitation(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	invitationController := do.MustInvoke[controller.InvitationController](injector)

	routes := route.Group("/api/invitation")
	{
		// Invitation
		routes.GET("/:id", middleware.Authenticate(jwtService), invitationController.GetInvitationByID)
		routes.GET("/event/:event_id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "ormawa"), invitationController.GetInvitationByEventID)
		routes.GET("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin", "ormawa"), invitationController.GetAllInvitations)
		routes.POST("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa"), invitationController.Create)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa"), invitationController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("ormawa"), invitationController.Delete)
	}
}
