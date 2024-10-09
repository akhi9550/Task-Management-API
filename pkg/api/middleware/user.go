package middleware

import (
	"net/http"
	"taskmanagementapi/pkg/helper"

	"github.com/gofiber/fiber/v2"
)

func UserAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenString := helper.GetTokenFromHeader(authHeader)

		if tokenString == "" {
			tokenString = c.Cookies("Authorization")
			if tokenString == "" {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"status":  http.StatusUnauthorized,
					"message": "Unauthorized",
					"data":    nil,
					"error":   "Authorization cookie not found",
				})
			}
		}

		userID, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  http.StatusUnauthorized,
				"message": "Invalid Token",
				"data":    nil,
				"error":   err.Error(),
			})
		}

		c.Locals("user_id", userID)
		c.Locals("email", userEmail)

		return c.Next()
	}
}
