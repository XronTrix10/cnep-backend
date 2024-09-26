package utils

import (
	"cnep-backend/source/models"
)

// ConvertToUserResponse converts a User model to a UserResponse model.
// This function can be used across different handlers to ensure consistent
// user data representation in API responses.
//
// Parameters:
//   - user: A pointer to the User model to be converted
//
// Returns:
//   - UserResponse: The converted UserResponse model
func ConvertToUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		Phone:             user.Phone,
		Address:           user.Address,
		Skills:            user.Skills,
		HelpedOthersCount: user.HelpedOthersCount,
		HelpReceivedCount: user.HelpReceivedCount,
		Rating:            user.Rating,
		Badges:            user.Badges,
		Designation:       user.Designation,
		IsVerified:        user.IsVerified,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}
