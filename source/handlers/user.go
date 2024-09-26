package handlers

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
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

func UpdateUserProfile(db *gorm.DB) fiber.Handler {
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

		// Define allowed fields for update
		allowedFields := map[string]bool{
			"name":                true,
			"phone":               true,
			"address":             true,
			"skills":              true,
			"designation":         true,
			"helped_others_count": true,
			"help_received_count": true,
			"rating":              true,
			"badges":              true,
		}

		// Filter out non-allowed fields and validate data
		filteredData := make(map[string]interface{})
		for key, value := range updateData {
			if allowedFields[key] {
				switch key {
				case "name", "phone", "address", "designation":
					if str, ok := value.(string); ok && str != "" {
						filteredData[key] = str
					}
				case "skills":
					if skills, ok := updateData["skills"].([]interface{}); ok {
						// Convert []interface{} to []string
						var skillsStr []string
						for _, skill := range skills {
							if str, ok := skill.(string); ok {
								skillsStr = append(skillsStr, str)
							} else {
								return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
									"error": "Invalid skill type",
								})
							}
						}
						filteredData[key] = pq.StringArray(skillsStr)
					}
				case "helped_others_count", "help_received_count":
					if count, ok := value.(float64); ok {
						filteredData[key] = uint(count)
					}
				case "rating":
					if rating, ok := value.(float64); ok {
						if rating > 5 || rating < 1 {
							return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"error": "Rating must be between 1 and 5",
							})
						}
						filteredData[key] = int(rating)
					}
				case "badges":
					if badges, ok := value.([]interface{}); ok {
						// Convert []interface{} to []int64
						intBadges := make([]int64, len(badges))
						for i, badge := range badges {
							if intBadge, ok := badge.(float64); ok {
								intBadges[i] = int64(intBadge)
							} else {
								return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
									"error": "Invalid badge type",
								})
							}
						}
						filteredData[key] = pq.Int64Array(intBadges)
					}
				}
			}
		}

		// Check if there are any valid fields to update
		if len(filteredData) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No valid fields to update",
			})
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching user profile",
			})
		}

		// Update only the allowed fields
		if err := db.Model(&user).Updates(filteredData).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error updating user profile",
			})
		}

		// Fetch the updated user from the database
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching updated user profile",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User profile updated successfully",
			"user":    utils.ConvertToUserResponse(&user),
		})
	}
}
