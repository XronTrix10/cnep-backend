package main

import (
	"log"

	"cnep-backend/pkg/lib"
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
	// Initialize SMTP
	lib.InitSMTP()

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
	port := cfg.ServerPort
	log.Fatal(app.Listen(":" + port))
}
