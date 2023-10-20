package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

// CreateClient creates and returns a new Redis client for a specific database (DB) number.
// It initializes the client with the Redis server address and password from environment variables.
//
// Parameters:
//   - dbNo: An integer representing the Redis database number to connect to.
//
// Returns:
//   - rdb: A pointer to a Redis client configured to connect to the specified database.
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}
