package main

import (
	"log"

	"cnep-backend/pkg/lib"
	"cnep-backend/source/config"
	"cnep-backend/source/database"
	"cnep-backend/source/routes"
	"cnep-backend/source/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize config
	cfg := config.New()

	// Initialize SMTP
	lib.InitSMTP()
	// Initialize Start time
	handlers.StatusInit()

	// Initialize database
	database.Connect(cfg)
	// Close database connection when the program exits
	defer database.Close()

	// Create Fiber app
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := cfg.ServerPort
	log.Printf("Server is starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
