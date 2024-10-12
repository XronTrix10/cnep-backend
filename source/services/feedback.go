package services

import (
	"cnep-backend/source/database"
	"cnep-backend/source/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

func AddFeedback(c *fiber.Ctx, senderID uint, receiverID uint, content string, rating uint8) error {
	var feedback models.Feedback

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if content == "" || senderID == receiverID || rating < 1 || rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	feedback.SenderID = senderID
	feedback.ReceiverID = receiverID
	feedback.Content = content
	feedback.Rating = rating

	if err := database.DB.Create(&feedback).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create feedback"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Feedback created successfully",
		"feedback": feedback,
	})
}

func GetFeedbackByUserID(c *fiber.Ctx, userID uint) error {
	var feedbacks []models.Feedback

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if err := database.DB.Where("receiver_id = ?", userID).Find(&feedbacks).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Feedback not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"feedbacks": feedbacks,
	})
}
