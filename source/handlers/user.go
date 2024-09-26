package handlers

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "cnep-backend/source/models"
)

func GetUserProfile(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get the user ID from the context (set by the AuthMiddleware)
        userID, ok := c.Locals("userID").(uint)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized: Invalid user ID",
            })
        }

        // Fetch the user from the database
        var user models.UserResponse
        if err := db.Table("users").First(&user, userID).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                    "error": "User not found",
                })
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Error fetching user profile",
            })
        }

        // Return the user profile
        return c.Status(fiber.StatusOK).JSON(user)
    }
}

