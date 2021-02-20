package lib

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads all the environment variables stored in a .env file
func LoadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}
}
