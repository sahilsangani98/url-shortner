package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sahilsangani98/url-shortner/routes"
)

// setupRoutes is a function that defines the routing for the application.
// It sets up URL paths and associates them with corresponding route handlers.
func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)    // Define a route to resolve shortened URLs
	app.Post("/api/v1", routes.ShortenURL) // Define an API endpoint to shorten URLs
}

func main() {
	// Load environment variables from a .env file (if available)
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New() // Create a new Fiber application instance

	app.Use(logger.New()) // Use the Fiber logger middleware to log HTTP requests

	setupRoutes(app) // Set up routing for the application

	// Start the Fiber application and listen on the port specified in the environment variables
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
