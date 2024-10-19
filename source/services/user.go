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
		"designation":         true,
		"rating":              true,
		"badges":              true,
		"topics":              true,
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
			case "rating":
				if rating, ok := value.(float64); ok {
					if rating > 5 || rating < 1 {
						return nil, fiber.NewError(fiber.StatusBadRequest, "Rating must be between 1 and 5")
					}
					filteredData[key] = float32(rating)
				}
			case "badges", "topics":
				if values, ok := value.([]interface{}); ok {
					// Convert []interface{} to []int64
					intValues := make([]int64, len(values))
					for i, badge := range values {
						if intBadge, ok := badge.(float64); ok {
							intValues[i] = int64(intBadge)
						} else {
							return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid badge or topic type")
						}
					}
					filteredData[key] = pq.Int64Array(intValues)
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
