package handlers

import (
	"cnep-backend/source/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

/*
GetUserProfile: Fetches the user profile for the authenticated user.
It takes the user ID from the context and fetches the user profile from the database.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the user profile.
*/
func GetUserProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid user ID",
			})
		}

		user, err := services.GetUserProfileByID(uint(userID))
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		// Return the user profile
		return c.Status(fiber.StatusOK).JSON(user)
	}
}

/*
GetUserProfileByID: Fetches the user profile for the specified user ID.
It takes the user ID from the URL parameter and fetches the user profile from the database.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the user profile.
*/
func GetUserProfileByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the URL parameter
		userIDParam := c.Params("id")

		// Validate that the user ID is an integer
		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
		}

		user, err := services.GetUserProfileByID(uint(userID))
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}

/*
UpdateUserProfile: Updates the user profile for the authenticated user.
It takes the user ID from the context, the updated data from the request body, and the database connection.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the updated user profile.
*/
func UpdateUserProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid user ID",
			})
		}

		var updateData map[string]interface{}
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		user, err := services.UpdateUserProfile(userID, updateData)
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User profile updated successfully",
			"user":    user,
		})
	}
}
