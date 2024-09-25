package handlers

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
The `RegenerateOTP` function is a handler function that regenerates the OTP for a user.
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
func RegenerateOTP(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if !utils.IsValidEmail(input.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		return regenerateOTP(c, db, &user)
	}
}

/*
The `regenerateOTP` function is a helper function that regenerates the OTP for a user.
It takes in the `c` context, the `db` database connection, and the `user` user object.
Here's a breakdown of what it does:

Steps:
 1. Generates a new OTP for the user.
 2. Updates the user's OTP and OTP expiry time in the database.
 3. Sends an email with the new OTP to the user.
 4. Returns a JSON response with a success message.

Returns:

	A JSON response with a success message if the OTP is regenerated, or an error message if the OTP cannot be regenerated.
*/
func regenerateOTP(c *fiber.Ctx, db *gorm.DB, user *models.User) error {
	otp := utils.GenerateOTP()
	otpExpiry := time.Now().Add(15 * time.Minute)

	user.OTP = otp
	user.OTPExpiry = otpExpiry

	if err := db.Save(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update OTP"})
	}

	if err := utils.SendOTPEmail(user.Email, user.Name, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send OTP email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New OTP sent. Please verify your email with the new OTP.",
	})
}

/*
The `VerifyOTP` function is a handler function that verifies the OTP sent to the user's email.
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
func VerifyOTP(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if !utils.IsValidEmail(input.Email) || len(input.OTP) != 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email or OTP format"})
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		if user.OTP != input.OTP || time.Now().After(user.OTPExpiry) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired OTP"})
		}

		user.IsVerified = true
		user.OTP = ""
		user.OTPExpiry = time.Time{}

		if err := db.Save(&user).Error; err != nil {
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
}
