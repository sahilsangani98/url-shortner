package routes

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sahilsangani98/url-shortner/database"
)

// ResolveURL resolves a short URL to its original full-length URL.
// It takes an HTTP request context (c *fiber.Ctx) and extracts the short URL
// from the request. The function queries a Redis database to find the original URL
// associated with the short URL. If a match is found, it increments a redirect counter
// and sends a 301 Moved Permanently redirect response to the client, taking them
// to the original URL.
//
// If the short URL is not found in the database, it returns a 404 Not Found response
// to indicate that the short URL does not exist.
//
// If there is a general database connectivity issue, it returns a 500 Internal Server
// Error response to indicate a problem connecting to the database.
//
// This function is part of a URL shortening service, and it is responsible for resolving
// short URLs and tracking the number of times they are accessed.
//
// Parameters:
//   - c: A *fiber.Ctx object representing the HTTP request context.
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	log.Printf("%s requested %s short url", c.IP(), url)

	rds0 := database.CreateClient(0)
	defer rds0.Close()

	value, err := rds0.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		log.Printf("%s requested short url - %s not found in database", c.IP(), url)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"Error": "Short URL not found in the database",
		})
	} else if err != nil {
		log.Printf("%s request couldn't severed as service is unable to connect to the database", c.IP())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Can not connect to database",
		})
	}

	// Increment the counter
	rds1 := database.CreateClient(1)
	defer rds1.Close()
	_ = rds1.Incr(database.Ctx, "counter")

	// Redirect to original URL
	return c.Redirect(value, 301)
}
