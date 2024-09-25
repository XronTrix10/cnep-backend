package main

import (
	"log"
	"os"

	"cnep-backend/source/config"
	"cnep-backend/source/database"
	"cnep-backend/source/routes"

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

	// Initialize database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create Fiber app
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}