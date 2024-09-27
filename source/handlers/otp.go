package handlers

import (
	"cnep-backend/source/services"

	"github.com/gofiber/fiber/v2"
)

/*
The `RegenerateOTP` function is a handler function that regenerates the OTP for a user.
Here's a breakdown of what it does:

Steps:
 1. Calls the OTPRegenerate function from the services package.
 2. Returns the JSON response with a success message if the OTP is regenerated, or an error message if the OTP cannot be regenerated.

Returns:

	A JSON response with a success message if the OTP is regenerated, or an error message if the OTP cannot be regenerated.
*/
func RegenerateOTP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.OTPRegenerate(c, input.Email)
	}
}

/*
The `VerifyOTP` function is a handler function that verifies the OTP sent to the user's email.
Here's a breakdown of what it does:

Steps:
 1. Retrieves the email and OTP parameters from the request body.
 2. Checks if the email and OTP parameters are valid email and OTP formats.
 3. Retrieves the user with the given email from the database.
 4. Checks if the user exists in the database.
 5. Checks if the OTP is valid for the user.
 6. If the OTP is valid, it updates the user's OTP and OTP expiry time in the database.
 7. If the OTP is not valid, it returns a JSON response with an error message.

Returns:

	A JSON response with a success message if the OTP is verified, or an error message if the OTP is invalid or expired.
*/
func VerifyOTP() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.ValidateOTP(c, input.OTP, input.Email)
	}
}
