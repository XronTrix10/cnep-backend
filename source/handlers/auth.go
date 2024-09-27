package handlers

import (
	"cnep-backend/source/services"

	"github.com/gofiber/fiber/v2"
)

/*
The `CheckEmailExistence` function is a handler function that checks the existence of a user
with a specific email address in the database.

Returns:

	A JSON response with a success message if the user exists, or an error message if the user does not exist.
*/
func CheckEmailExistence() fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email")

		return services.CheckEmailExistence(c, email)
	}
}

/*
The `Authentication` function is a handler function that registers or logs in a user with the system.

Returns:

	A JSON response with a success message if the user is registered, or an error message if the user cannot be registered.
*/
func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			IsNew    bool   `json:"is_new"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if input.IsNew {
			return services.RegisterService(c, input.Email, input.Password)
		}

		return services.LoginService(c, input.Email, input.Password)
	}
}
