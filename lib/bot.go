package lib

import (
	"github.com/go-redis/redis/v8"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot is a structure that holds the Telegram API and other assets and is unique
type Bot struct {
	// The Telegram API that is merged into our structure
	*telegram.BotAPI

	// A list of the slash commands available to the user
	Commands map[string]Command
	// The Redis client the bot will connect to
	Cache *redis.Client
}

// Help sends a formatted help message in the Telegram chat
func (bot *Bot) Help(chatID int64, command string) (err error) {
	help := telegram.NewMessage(chatID, bot.Commands[command].Help())
	help.ParseMode = "Markdown"
	_, err = bot.Send(help)
	return
}

// Error sends a formatted error message in the Telegram chat
func (bot *Bot) Error(chatID int64, message string) (err error) {
	msg := telegram.NewMessage(chatID, "❌ "+message)
	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)
	return
}
