package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/constants"
	"github.com/miraicantsleep/myits-event-be/controller"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/samber/do"
)

func Department(route *gin.Engine, injector *do.Injector) {
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)
	departmentController := do.MustInvoke[controller.DepartmentController](injector)

	routes := route.Group("/api/department")
	{
		// Department
		routes.POST("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin"), departmentController.Create)
		routes.GET("/", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin"), departmentController.GetAllDepartment)
		routes.GET("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin"), departmentController.GetDepartmentByID)
		routes.PATCH("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin"), departmentController.Update)
		routes.DELETE("/:id", middleware.Authenticate(jwtService), middleware.RoleMiddleware("admin"), departmentController.Delete)
	}
}
