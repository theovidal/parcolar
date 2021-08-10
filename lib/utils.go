package lib

import (
	"net/url"
	"regexp"

	"github.com/joho/godotenv"
)

// LoadEnv loads all the environment variables stored in a .env file
func LoadEnv(path string) {
	_ = godotenv.Load(path)
}

// EncodeURL parses an URL with its query parameters
func EncodeURL(input string, params map[string]string) string {
	url, _ := url.Parse(input)
	query := url.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	url.RawQuery = query.Encode()
	return url.String()
}

// ParseTelegramMessage escapes all the characters required to print MarkdownV2 content
func ParseTelegramMessage(input string) (output string) {
	escape := regexp.MustCompile("[_\\*\\[\\]\\(\\)~>#\\+-=\\|{}\\.!]+")
	for _, char := range []rune(input) {
		if escape.Match([]byte(string(char))) {
			output += "\\"
		}
		output += string(char)
	}
	return
}

// Contains check if a specific slice contains a string
func Contains(slice []string, text string) bool {
	for _, item := range slice {
		if item == text {
			return true
		}
	}

	return false
}
