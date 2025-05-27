package main

import (
	"log"
	"os"

	"github.com/miraicantsleep/myits-event-be/command"
	"github.com/miraicantsleep/myits-event-be/middleware"
	"github.com/miraicantsleep/myits-event-be/provider"
	"github.com/miraicantsleep/myits-event-be/routes"
	"github.com/samber/do"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
)

func args(injector *do.Injector) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(injector)
		return flag
	}

	return true
}

func run(server *gin.Engine) {
	server.Static("/assets", "./assets")

	if os.Getenv("IS_LOGGER") == "true" {
		routes.LoggerRoute(server)
	}

	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "0.0.0.0:" + port
	} else {
		serve = ":" + port
	}

	myFigure := figure.NewColorFigure("Caknoo", "", "green", true)
	myFigure.Print()

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func main() {
	var (
		injector = do.New()
	)

	provider.RegisterDependencies(injector)

	if !args(injector) {
		return
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// add ping
	server.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// routes
	routes.RegisterRoutes(server, injector)

	run(server)
}
