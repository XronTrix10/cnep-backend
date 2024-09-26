package middleware

import (
	"strings"

	"cnep-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the Authorization header is empty or doesn't start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing or invalid Authorization header",
			})
		}

		// Extract the token from the Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token",
			})
		}

		// Add the user ID to the context for use in subsequent handlers
		c.Locals("userID", userID)
		return c.Next()
	}
}
