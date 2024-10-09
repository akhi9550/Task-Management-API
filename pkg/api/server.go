package http

import (
	"log"
	"taskmanagementapi/pkg/api/handlers"
	"taskmanagementapi/pkg/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ServerHTTP struct {
	app *fiber.App
}

func NewServerHTTP(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler) *ServerHTTP {
	app := fiber.New(fiber.Config{})
	app.Use(logger.New())
	routes.UserRoutes(app.Group("/user"), userHandler)
	routes.TaskRoutes(app.Group("/tasks"), taskHandler)
	return &ServerHTTP{app: app}
}

func (sh *ServerHTTP) Start(infoLog *log.Logger, errorLog *log.Logger) {
	infoLog.Printf("starting server on :3000")
	err := sh.app.Listen(":3000")
	if err != nil {
		errorLog.Fatal(err)
	}
}
