package lib

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}
}
