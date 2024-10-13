package handlers

import (
	"cnep-backend/source/services"
	"cnep-backend/pkg/template"
	"github.com/gofiber/fiber/v2"
)

/*
The `ChangePassword` function is a handler function that allows users to change their password.
It takes the user ID from the context, the updated data from the request body, and the database connection.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the updated user profile.
*/
func ChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		var input struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.ChangePassword(c, userID, input.OldPassword, input.NewPassword)
	}
}
