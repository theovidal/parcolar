package lib

import (
	"github.com/joho/godotenv"
)

// LoadEnv loads all the environment variables stored in a .env file
func LoadEnv(path string) {
	_ = godotenv.Load(path)
}
