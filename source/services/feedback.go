package services

import (
	"cnep-backend/pkg/consts"
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
	var feedbacks []models.FeedbackSender

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	err := database.DB.Table(consts.FEEDBACK_TABLE).
		Select("feedbacks.id as feedback_id, feedbacks.content, feedbacks.rating, feedbacks.sender_id, users.name as sender_name, users.email as sender_email, users.rating as sender_rating, feedbacks.created_at, feedbacks.updated_at").
		Joins("JOIN users ON feedbacks.sender_id = users.id").
		Where("feedbacks.receiver_id = ?", userID).
		Scan(&feedbacks).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Feedback not found"})
	}

	var nestedFeedbacks []models.FeedbackWithSender

	for _, feedback := range feedbacks {
		feedback := models.FeedbackWithSender{
			FeedbackID: feedback.FeedbackID,
			Content:    feedback.Content,
			Rating:     feedback.Rating,
			Sender: models.SenderInfo{
				ID:     feedback.SenderID,
				Name:   feedback.SenderName,
				Email:  feedback.SenderEmail,
				Rating: feedback.SenderRating,
			},
			CreatedAt: feedback.CreatedAt,
			UpdatedAt: feedback.UpdatedAt,
		}
		nestedFeedbacks = append(nestedFeedbacks, feedback)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"feedbacks": nestedFeedbacks,
	})
}
