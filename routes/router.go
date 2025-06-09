package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	User(server, injector)
	Department(server, injector)
	Event(server, injector)
	Room(server, injector)
	Invitation(server, injector)
}
