package lib

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv(path string) {
	if err := godotenv.Load(path); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}
}
