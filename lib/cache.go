package lib

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// Cache is the main Redis cache of the app
var Cache *redis.Client

// OpenCache creates the Redis client and connects to the cache
func OpenCache() {
	db, err := strconv.Atoi(os.Getenv("CACHE_DB"))
	if err != nil {
		log.Panicln("‼ Error creating the cache: you must provide an integer for DB's number")
	}
	Cache = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_ADDRESS") + ":" + os.Getenv("CACHE_PORT"),
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       db,
	})
}
