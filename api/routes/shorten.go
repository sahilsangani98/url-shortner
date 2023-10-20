package routes

import (
	// "fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sahilsangani98/url-shortner/database"
	"github.com/sahilsangani98/url-shortner/helpers"
)

// Define the structure for incoming requests.
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// Define the structure for responses.
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL handles incoming HTTP requests.
func ShortenURL(c *fiber.Ctx) error {
	// Check for incomming body request and validate against pre-defined `request`` structure
	body := new(request) // Creates zero-value filled structure for new request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Can not parse JSON",
		})
	}

	// Log that the request has been received.
	log.Printf("Received request from IP %s for URL %s, CustomShort %s, & Expiry %s", c.IP(), body.URL, body.CustomShort, body.Expiry)

	// Implement rate limiting:
	// - Check if the user's IP is already in the database.
	// - If yes, decrement the remaining calls by one.
	// - If no, add the IP to the database with an expiry of 30 minutes.
	//   This allows the user to send 10 requests every 30 minutes.

	rds1 := database.CreateClient(1)
	defer rds1.Close()

	val, err := rds1.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = rds1.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := rds1.TTL(database.Ctx, c.IP()).Result()
			log.Printf("Rate limit exceeded for %s", c.IP())
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"Error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// Validate provided URL
	if !govalidator.IsURL(body.URL) {
		log.Printf("Accessed invalid URL by %s", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Invalid URL",
		})
	}

	// Validate the domain
	if !helpers.RemoveDomainError(body.URL) {
		log.Printf("Accessed invalid domain for %s", c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": "Invalid domain",
		})
	}

	// Enforce https and SSL on the URL
	body.URL = helpers.EnforceHTTP(body.URL)

	// Check if the user has provided any custom short URLs.
	// - If yes, proceed.
	// - If no, create a new short URL using the first 6 digits of UUID.
	// - No collision checks are performed in this code.
	//   You can implement collision checks as needed.
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	rds0 := database.CreateClient(0)
	defer rds0.Close()

	val, _ = rds0.Get(database.Ctx, id).Result()

	// Check if provided short URL is already in use
	if val != "" {
		log.Printf("%s provided in already use URL", c.IP())
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"Error": "Provided short URL already in use. Please provide some other short URL",
		})
	}

	if body.Expiry == 0 {
		expiry, _ := strconv.Atoi(os.Getenv("URL_RETENTION_TIME")) // default expiry of 24 hours
		body.Expiry = time.Duration(expiry)
	}

	err = rds0.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		log.Printf("%s request couldn't severed as service is unable to connect to the server", c.IP())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Unable to connect to the server",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	rds1.Decr(database.Ctx, c.IP())

	val, _ = rds1.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := rds1.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
