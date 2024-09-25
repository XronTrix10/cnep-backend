package handlers

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CheckEmailExistence checks if a user with the given email exists
func CheckEmailExistence(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email")

		if email == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email parameter is required"})
		}

		if !utils.IsValidEmail(email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		}

		var user models.User
		result := db.Where("email = ?", email).First(&user)

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
}

func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		if !utils.IsValidEmail(user.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		}

		if !utils.IsValidPassword(user.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password does not meet complexity requirements"})
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
		}
		user.Password = hashedPassword

		// Generate OTP
		otp := utils.GenerateOTP()                    // Implement this function in your utilss package
		otpExpiry := time.Now().Add(15 * time.Minute) // Expiry time for OTP is 15 minutes

		user.OTP = otp
		user.OTPExpiry = otpExpiry
		user.IsVerified = false

		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
		}

		// TODO: Which feilds to take input fro user at registration ?

		// Send OTP via email
		if err := utils.SendOTPEmail(user.Email, "John Doe", otp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send OTP email"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered. Please verify your email with the OTP sent.",
		})
	}
}

func VerifyOTP(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
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

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		if !user.IsVerified {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Email not verified",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
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
			"user":  user,
		})
	}
}

func RegenerateOTP(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		var user models.User
		if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		return regenerateOTP(c, db, &user)
	}
}

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
