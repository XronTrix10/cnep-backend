package routes

import (
	"cnep-backend/source/handlers"
	"cnep-backend/source/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// use logger middleware
	app.Use(middleware.Logger())

	app.Get("/", handlers.Status())
	
	// Public routes
	app.Get("/api/auth/users", handlers.CheckEmailExistence())
	app.Post("/api/auth/continue", handlers.Authentication())
	app.Post("/api/otp/generate", handlers.RegenerateOTP())
	app.Post("/api/otp/verify", handlers.VerifyOTP())

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())

	// ===================================================================

	usersApi := api.Group("/users")

	// User routes
	usersApi.Get("/profile", handlers.GetUserProfile())
	usersApi.Get("/profile/:id", handlers.GetUserProfileByID())
	usersApi.Put("/profile", handlers.UpdateUserProfile())
	
	// Sensitive routes
	usersApi.Post("/password/change", handlers.ChangePassword())

	// Feedback routes
	usersApi.Post("/feedback", handlers.AddFeedback())
	usersApi.Get("/feedback", handlers.GetFeedback())
	usersApi.Get("/feedback/:id", handlers.GetFeedbackByID())

	// Partner routes
	usersApi.Get("/partner", handlers.GetPartners())
	usersApi.Post("/partner", handlers.AddPartner())
	usersApi.Put("/partner/:id", handlers.UpdatePartnerStatus())
	usersApi.Get("/partner/pending", handlers.GetPendingPartners())
	usersApi.Delete("/partner/:id", handlers.CancelPartnerRequest())

	// // Post routes
	// api.Get("/posts", handlers.GetPosts(db))
	// api.Post("/posts", handlers.CreatePost(db))
	// api.Get("/posts/:id", handlers.GetPost(db))
	// api.Put("/posts/:id", handlers.UpdatePost(db))
	// api.Delete("/posts/:id", handlers.DeletePost(db))

	// // Comment routes
	// api.Get("/posts/:id/comments", handlers.GetComments(db))
	// api.Post("/posts/:id/comments", handlers.CreateComment(db))

	// // Conversation routes
	// api.Get("/conversations", handlers.GetConversations(db))
	// api.Post("/conversations", handlers.CreateConversation(db))

	// // WebSocket route for chat
	// app.Get("/ws", middleware.AuthMiddleware(), handlers.HandleWebSocket())
}
