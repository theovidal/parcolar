package lib

import (
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
)

var Cache *redis.Client

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
