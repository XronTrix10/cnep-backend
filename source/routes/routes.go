package routes

import (
	"cnep-backend/source/handlers"
	"cnep-backend/source/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// use logger middleware
	app.Use(middleware.Logger())
	
	// Public routes
	app.Get("/api/auth/users", handlers.CheckEmailExistence())
	app.Post("/api/auth/continue", handlers.Authentication())
	app.Post("/api/otp/generate", handlers.RegenerateOTP())
	app.Post("/api/otp/verify", handlers.VerifyOTP())

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())

	// User routes
	api.Get("/users/profile", handlers.GetUserProfile())
	api.Get("/users/profile/:id", handlers.GetUserProfileByID())
	api.Put("/users/profile", handlers.UpdateUserProfile())

	// // Post routes
	// api.Get("/posts", handlers.GetPosts(db))
	// api.Post("/posts", handlers.CreatePost(db))
	// api.Get("/posts/:id", handlers.GetPost(db))
	// api.Put("/posts/:id", handlers.UpdatePost(db))
	// api.Delete("/posts/:id", handlers.DeletePost(db))

	// // Comment routes
	// api.Get("/posts/:id/comments", handlers.GetComments(db))
	// api.Post("/posts/:id/comments", handlers.CreateComment(db))

	// // Connection routes
	// api.Post("/connections", handlers.CreateConnection(db))
	// api.Put("/connections/:id", handlers.UpdateConnection(db))

	// // Conversation routes
	// api.Get("/conversations", handlers.GetConversations(db))
	// api.Post("/conversations", handlers.CreateConversation(db))

	// // WebSocket route for chat
	// app.Get("/ws", middleware.AuthMiddleware(), handlers.HandleWebSocket())
}
