package routes

import (
	"cnep-backend/source/handlers"
	"cnep-backend/source/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Public routes
	app.Post("/api/register", handlers.Register(db))
	app.Post("/api/login", handlers.Login(db))

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware())

	// User routes
	api.Get("/users", handlers.GetUsers(db))
	api.Get("/users/:id", handlers.GetUser(db))
	api.Put("/users/:id", handlers.UpdateUser(db))

	// Post routes
	api.Get("/posts", handlers.GetPosts(db))
	api.Post("/posts", handlers.CreatePost(db))
	api.Get("/posts/:id", handlers.GetPost(db))
	api.Put("/posts/:id", handlers.UpdatePost(db))
	api.Delete("/posts/:id", handlers.DeletePost(db))

	// Comment routes
	api.Get("/posts/:id/comments", handlers.GetComments(db))
	api.Post("/posts/:id/comments", handlers.CreateComment(db))

	// Connection routes
	api.Post("/connections", handlers.CreateConnection(db))
	api.Put("/connections/:id", handlers.UpdateConnection(db))

	// Conversation routes
	api.Get("/conversations", handlers.GetConversations(db))
	api.Post("/conversations", handlers.CreateConversation(db))

	// WebSocket route for chat
	app.Get("/ws", middleware.AuthMiddleware(), handlers.HandleWebSocket())
}
