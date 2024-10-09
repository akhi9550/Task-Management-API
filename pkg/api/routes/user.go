package routes

import (
	"taskmanagementapi/pkg/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router, userHandler *handlers.UserHandler) {
	app.Post("/signup", userHandler.UserSignUp)
	app.Post("/signin", userHandler.UserSignIn)
	app.Post("/signout", userHandler.UserSignOut)

}
