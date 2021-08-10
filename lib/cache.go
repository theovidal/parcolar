package lib

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// OpenCache creates the Redis client and connects to the cache
func OpenCache() (cache *redis.Client) {
	db, err := strconv.Atoi(os.Getenv("CACHE_DB"))
	if err != nil {
		Fatal("Error creating the cache: you must provide an integer for DB's number")
	}
	cache = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_ADDRESS") + ":" + os.Getenv("CACHE_PORT"),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       db,
	})

	if err = cache.Ping(context.Background()).Err(); err != nil {
		Fatal("Error connecting to the cache: %s", err)
	}
	LogInfo("Opened Redis cache on port %s", os.Getenv("CACHE_PORT"))
	return
}
