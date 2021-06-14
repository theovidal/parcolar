package lib

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads all the environment variables stored in a .env file
func LoadEnv(path string) {
	var message string
	if err := godotenv.Load(path); err != nil {
		message = "No .env file at the root - Ignoring"
	} else {
		message = "Loaded environment variables from .env file"
	}
	log.Println(Cyan.Sprintf("🏷 " + message))
}
