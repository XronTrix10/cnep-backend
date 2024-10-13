package handlers

import (
	"cnep-backend/source/services"
	"cnep-backend/pkg/template"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

/*
The `AddFeedback` function is a handler function that allows users to submit feedback.
It takes the user IDs from the context, the feedback content, and the rating from the request body.
If the user is not found or any other error occurs, it returns an appropriate error message.
The function returns a JSON response with the feedback.
*/
func AddFeedback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			UserID  uint   `json:"user_id"`
			Content string `json:"content"`
			Rating  uint8  `json:"rating"`
		}

		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.AddFeedback(c, userID, input.UserID, input.Content, input.Rating)
	}
}

func GetFeedback() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		return services.GetFeedbackByUserID(c, userID)
	}
}

func GetFeedbackByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		UserID := c.Params("id")

		if UserID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID not provided in query"})
		}

		IntUserID, err := strconv.Atoi(UserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
		}
		UintUserID := uint(IntUserID)
		return services.GetFeedbackByUserID(c, UintUserID)
	}
}
