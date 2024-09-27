package services

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/database"
	"cnep-backend/source/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*
The OTPRegenerate function is a handler function that regenerates the OTP for a user.
Here's a breakdown of what it does:

Steps:
 1. Retrieves the email parameter from the request body.
 2. Checks if the email parameter is empty or not.
 3. Checks if the email parameter is a valid email format.
 4. Retrieves the user with the given email from the database.
 5. Checks if the user exists in the database.
 6. Generates a new OTP for the user.
 7. Updates the user's OTP and OTP expiry time in the database.
 8. Sends an email with the new OTP to the user.
 9. Returns a JSON response with a success message.

Returns:

	A JSON response with a success message if the OTP is regenerated, or an error message if the OTP cannot be regenerated.
*/
func OTPRegenerate(c *fiber.Ctx, email string) error {
	var user models.User

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidEmail(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(15 * time.Minute)

	user.OTP = otp
	user.OTPExpiry = otpExpiry

	if err := database.DB.Save(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update OTP"})
	}

	if err := utils.SendOTPEmail(user.Email, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send OTP email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New OTP sent. Please verify your email with the new OTP.",
	})
}

/*
The ValidateOTP function is a handler function that verifies the OTP sent to the user's email.
Here's a breakdown of what it does:

Steps:

 1. Retrieves the email and OTP parameters from the request body.
 2. Checks if the email and OTP parameters are valid email and OTP formats.
 3. Retrieves the user with the given email from the database.
 4. Checks if the user exists in the database.
 5. Checks if the OTP is valid for the user.
 6. If the OTP is valid, it updates the user's OTP and OTP expiry time in the database.
 7. If the OTP is not valid, it returns a JSON response with an error message.

Returns:

	A JSON response with a success message if the OTP is verified, or an error message if the OTP is invalid or expired.
*/
func ValidateOTP(c *fiber.Ctx, otp string, email string) error {
	var user models.User

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidEmail(email) || len(otp) != 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or OTP format"})
	}

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if user.OTP != otp || time.Now().After(user.OTPExpiry) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	user.IsVerified = true
	user.OTP = ""
	user.OTPExpiry = time.Time{}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not verify user"})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{
		"message": "Email verified successfully",
		"token":   token,
	})
}
