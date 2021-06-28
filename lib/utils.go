package lib

import (
	"regexp"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Error sends a formatted error message in the Telegram chat
func Error(bot *telegram.BotAPI, update *telegram.Update, message string) (err error) {
	msg := telegram.NewMessage(update.Message.Chat.ID, "❌ "+message)
	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)
	return
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
