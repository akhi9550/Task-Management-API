package routes

import (
	"taskmanagementapi/pkg/api/handlers"
	"taskmanagementapi/pkg/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(app fiber.Router, taskHandler *handlers.TaskHandler) {
	app.Use(middleware.UserAuthMiddleware())
	{
		app.Post("", taskHandler.CreateTask)
		app.Get("", taskHandler.GetTasks)
		app.Get("/:id", taskHandler.GetTask)
		app.Put("/:id", taskHandler.UpdateTask)
		app.Delete("/:id", taskHandler.DeleteTask)
	}
}
