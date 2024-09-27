package services

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/database"
	"cnep-backend/source/models"
	"github.com/gofiber/fiber/v2"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

// Checks if the email exists in the database
func CheckEmailExistence(c *fiber.Ctx, email string) error {
	var user models.User

	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidEmail(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	result := database.DB.Where("email = ?", email).First(&user)

	if result.Error == nil {
		// User exists
		return c.JSON(fiber.Map{"exists": true})
	} else if result.Error == gorm.ErrRecordNotFound {
		// User does not exist
		return c.JSON(fiber.Map{"exists": false})
	} else {
		// Database error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
}

// Registers a new user in the database
func RegisterService(c *fiber.Ctx, email, password string) error {
	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidEmail(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
	}

	if !utils.IsValidPassword(password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password does not meet complexity requirements"})
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	var user models.User
	user.Password = hashedPassword
	user.Email = email

	// Generate OTP
	otp := utils.GenerateOTP()                    // Implement this function in your utilss package
	otpExpiry := time.Now().Add(15 * time.Minute) // Expiry time for OTP is 15 minutes

	user.OTP = otp
	user.OTPExpiry = otpExpiry
	user.IsVerified = false

	if err := database.DB.Create(&user).Error; err != nil {
		if utils.IsDuplicateEntryError(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(user.Email, otp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send OTP email"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered. Please verify your email with the OTP sent.",
	})
}

// Logs in a user with the provided email and password
func LoginService(c *fiber.Ctx, email, password string) error {
	// Ensure database connection is established
	if database.DB == nil {
		log.Panic("Database not connected")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database not connected"})
	}

	if !utils.IsValidEmail(email) || !utils.IsValidPassword(password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if !user.IsVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email not verified",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user":  utils.ConvertToUserResponse(&user),
	})
}
