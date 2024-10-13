package template

import "github.com/gofiber/fiber/v2"

func Unauthenticated(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized: Missing or invalid Authorization header",
	})
}