package handlers

import (
	"cnep-backend/pkg/utils"
	"cnep-backend/source/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

/*
The `CheckEmailExistence` function is a handler function that checks the existence of a user
with a specific email address in the database. Here's a breakdown of what it does:

Steps:
 1. Retrieves the email parameter from the request body.
 2. Checks if the email parameter is empty or not.
 3. Checks if the email parameter is a valid email format.
 4. Retrieves the user with the given email from the database.
 5. Checks if the user exists in the database.
 6. If the user exists, it returns a JSON response with a success message.
 7. If the user does not exist, it returns a JSON response with an error message.

Returns:

	A JSON response with a success message if the user exists, or an error message if the user does not exist.
*/
func CheckEmailExistence(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email")

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

/*
The `Register` function is a handler function that registers a new user with the system.
Here's a breakdown of what it does:

Steps:
 1. Retrieves the user object from the request body.
 2. Checks if the user object is valid.
 3. Checks if the email is a valid email format.
 4. Checks if the email is already in use.
 5. Checks if the password is a valid password format.
 6. Hashes the password.
 7. Generates a random OTP for the user.
 8. Sends an email with the OTP to the user.
 9. Returns a JSON response with a success message.

Returns:

	A JSON response with a success message if the user is registered, or an error message if the user cannot be registered.
*/
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
			if utils.IsDuplicateEntryError(err) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
			}
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

/*
The `Login` function is a handler function that logs in a user with their email and password.
Here's a breakdown of what it does:

Steps:
 1. Retrieves the email and password parameters from the request body.
 2. Checks if the email and password parameters are empty or not.
 3. Checks if the email and password parameters are valid email and password formats.
 4. Retrieves the user with the given email from the database.
 5. Checks if the user exists in the database.
 6. Checks if the password is valid for the user.
 7. Generates a JWT token for the user.
 8. Returns a JSON response with the token.

Returns:

	A JSON response with the token if the user is logged in, or an error message if the user cannot be logged in.
*/
func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Request Body",
			})
		}

		if !utils.IsValidEmail(input.Email) || !utils.IsValidPassword(input.Password) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email or password",
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
			"user":  utils.ConvertToUserResponse(&user),
		})
	}
}
