package handlers

import (
	"cnep-backend/pkg/template"
	"cnep-backend/source/services"

	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddPartner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			UserID uint `json:"user_id"`
		}

		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.AddPartner(c, userID, input.UserID)
	}
}

func UpdatePartnerStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the URL parameter
		PartnerIDParam := c.Params("id")

		// Validate that the user ID is an integer
		partnerID, err := strconv.Atoi(PartnerIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
		}

		var input struct {
			Accept bool `json:"accept"`
		}

		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		return services.UpdatePartnerStatus(c, userID, uint(partnerID), input.Accept)
	}
}

func CancelPartnerRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the URL parameter
		PartnerIDParam := c.Params("id")

		// Validate that the user ID is an integer
		partnerID, err := strconv.Atoi(PartnerIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
		}

		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		return services.CancelPartnerRequest(c, userID, uint(partnerID))
	}
}

func GetPartners() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		return services.GetPartners(c, userID)
	}
}

func GetPendingPartners() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the context (set by the AuthMiddleware)
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return template.Unauthenticated(c)
		}

		return services.GetPendingPartners(c, userID)
	}
}
