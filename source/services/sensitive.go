package services

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/database"
	"cnep-backend/source/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

/*
The ChangePassword function is a handler function that allows users to change their password.
It takes the user ID, the old password, and the new password as parameters.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the updated user profile.
*/
func ChangePassword(c *fiber.Ctx, userId uint, oldPassword string, newPassword string) error {
	var user models.User

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidPassword(oldPassword) || !utils.IsValidPassword(newPassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password format"})
	}

	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if !utils.CheckPasswordHash(user.Password, oldPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect old password"})
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user.Password = hashedPassword

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password successfully updated",
	})
}
