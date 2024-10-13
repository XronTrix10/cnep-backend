package services

import (
	"cnep-backend/pkg/consts"
	"cnep-backend/source/database"
	"cnep-backend/source/models"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
)

/*
The GetUserProfileByID function fetches the user profile for the specified user ID from the database.
It takes the user ID as a parameter and returns a pointer to a UserResponse struct.
If the user is not found or any other error occurs, it returns an appropriate error message.
*/
func GetUserProfileByID(userId uint) (*models.UserResponse, error) {
	var user models.UserResponse
	if database.DB == nil {
		log.Fatal("Database not connected")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Database not connected")
	}

	if err := database.DB.Table(consts.USERS_TABLE).First(&user, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user profile")
	}
	return &user, nil
}

/*
The UpdateUserProfile function updates the user profile for the specified user ID in the database.
It takes the user ID, the updated data, and the database connection as parameters.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the updated user profile.
*/
func UpdateUserProfile(userId uint, updateData map[string]interface{}) (*models.UserResponse, error) {
	var user models.User
	var userResponse models.UserResponse
	if database.DB == nil {
		log.Fatal("Database not connected")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Database not connected")
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
							return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid skill type")
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
						return nil, fiber.NewError(fiber.StatusBadRequest, "Rating must be between 1 and 5")
					}
					filteredData[key] = float32(rating)
				}
			case "badges":
				if badges, ok := value.([]interface{}); ok {
					// Convert []interface{} to []int64
					intBadges := make([]int64, len(badges))
					for i, badge := range badges {
						if intBadge, ok := badge.(float64); ok {
							intBadges[i] = int64(intBadge)
						} else {
							return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid badge type")
						}
					}
					filteredData[key] = pq.Int64Array(intBadges)
				}
			}
		}
	}

	// Check if there are any valid fields to update
	if len(filteredData) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No valid fields to update")
	}

	if err := database.DB.Table(consts.USERS_TABLE).First(&user, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching user profile")
	}

	// Update only the allowed fields
	if err := database.DB.Table(consts.USERS_TABLE).Model(&user).
		Updates(filteredData).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error updating user profile")
	}

	// Fetch the updated user from the database
	if err := database.DB.Table(consts.USERS_TABLE).
		First(&userResponse, userId).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error fetching updated user profile")
	}

	return &userResponse, nil
}
