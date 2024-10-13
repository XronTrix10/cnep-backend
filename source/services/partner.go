package services

import (
	"cnep-backend/source/database"
	"cnep-backend/source/models"

	"github.com/gofiber/fiber/v2"
	"log"
)

const (
	PARTNER_STATUS_PENDING  = "pending"
	PARTNER_STATUS_ACCEPTED = "accepted"
	PARTNER_STATUS_DECLINED = "declined"
	TABLE_NAME              = "partners"
)

func AddPartner(c *fiber.Ctx, senderID, receiverID uint) error {
	var partner models.Partner

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if senderID == receiverID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
	}

	partner.SenderID = senderID
	partner.ReceiverID = receiverID
	partner.Status = PARTNER_STATUS_PENDING

	if err := database.DB.Table(TABLE_NAME).Create(&partner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create partner"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "ok",
		"message": "Partner Request Created",
	})
}

func UpdatePartnerStatus(c *fiber.Ctx, userID, partnerID uint, accepted bool) error {
	var partner models.Partner

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if err := database.DB.Table(TABLE_NAME).Where("receiver_id = ? AND sender_id = ?", userID, partnerID).
		First(&partner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User or Partner not found"})
	}

	if accepted {
		partner.Status = PARTNER_STATUS_ACCEPTED
	} else {
		partner.Status = PARTNER_STATUS_DECLINED
	}

	if err := database.DB.Table(TABLE_NAME).Save(&partner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update partner status"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Partner status updated",
	})
}

func CancelPartnerRequest(c *fiber.Ctx, userID, partnerID uint) error {
	var partner models.Partner

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if err := database.DB.Table(TABLE_NAME).Where("receiver_id = ? AND sender_id = ?", partnerID, userID).
		First(&partner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Partner request does not exist"})
	}

	if err := database.DB.Table(TABLE_NAME).Delete(&partner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not cancel partner request"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "Partner request cancelled",
	})
}

func GetPartners(c *fiber.Ctx, userID uint) error {
	var partners []models.Partner
	var ids []uint
	var users []models.UserResponse

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if err := database.DB.Table(TABLE_NAME).
		Where("receiver_id = ? OR sender_id = ? AND status = ?", userID, userID, PARTNER_STATUS_ACCEPTED).
		Find(&partners).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Partner request not found"})
	}

	for _, partner := range partners {
		if userID == partner.SenderID {
			ids = append(ids, partner.ReceiverID)
		} else {
			ids = append(ids, partner.SenderID)
		}
	}

	// If no partners found, return an empty array
	if len(ids) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":   "empty array",
			"partners": users,
		})
	}

	// Retrieve user information for the collected IDs
	if err := database.DB.Table("users").Where("id IN ?", ids).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user information"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "ok",
		"partners": users,
	})
}

func GetPendingPartners(c *fiber.Ctx, userID uint) error {
	var partners []models.Partner
	var ids []uint
	var users []models.UserResponse

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if err := database.DB.Table(TABLE_NAME).
		Where("receiver_id = ? AND status = ?", userID, PARTNER_STATUS_PENDING).
		Find(&partners).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":   "No pending partners found",
			"partners": partners,
		})
	}

	for _, partner := range partners {
		ids = append(ids, partner.SenderID)
	}

	// If no partners found, return an empty array
	if len(ids) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":   "empty array",
			"partners": users,
		})
	}

	// Retrieve user information for the collected IDs
	if err := database.DB.Table("users").Where("id IN ?", ids).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve user information"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "ok",
		"partners": users,
	})
}
